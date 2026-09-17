package websocket

import (
	"backend/entities"
	"backend/entities/packet"
	"backend/metrics"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Configuration constants for connection management
const (
	// Maximum time to wait for a write operation
	writeWait = 10 * time.Second

	// Maximum time to wait for a pong response
	pongWait = 60 * time.Second

	// Interval for sending ping messages (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Clients only receive server-authoritative events; data frames are ignored.
	maxMessageSize = 1024
)

// HandleConnection handles new incoming HTTP connections and upgrades them to WebSockets
func (cs *chatsocket) RunWebsocket(conn *websocket.Conn, connUserID string) {
	// INITIALIZE CLIENT
	client := &entities.Client{
		Hub:  cs.hub,
		Conn: conn,
		Send: make(chan []byte, 256), // Buffered to prevent blocking other goroutines
		ID:   connUserID,
	}

	// Register with Hub
	cs.hub.Register <- client
	connectionID := uuid.NewString()
	if cs.registry != nil {
		if err := cs.registry.RegisterConnection(cs.hub.Ctx, connUserID, connectionID, cs.instance); err != nil {
			slog.Warn("failed to register websocket connection", "user_id", connUserID, "error", err)
		}
		closed := make(chan struct{})
		go cs.refreshConnection(connUserID, connectionID, closed)
		go func() {
			defer close(closed)
			readPump(client)
		}()
	} else {
		go readPump(client)
	}

	// RUN READ AND WRITE PUMPS IN SEPARATE GOROUTINES
	// This separates reading from writing, resolving concurrent-write safety issues.
	go writePump(client)
}

func (cs *chatsocket) refreshConnection(userID, connectionID string, closed <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-closed:
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			if err := cs.registry.UnregisterConnection(ctx, userID, connectionID, cs.instance); err != nil {
				slog.Warn("failed to unregister websocket connection", "user_id", userID, "error", err)
			}
			cancel()
			return
		case <-cs.hub.Ctx.Done():
			return
		case <-ticker.C:
			if err := cs.registry.RefreshConnection(cs.hub.Ctx, userID, connectionID); err != nil {
				slog.Warn("failed to refresh websocket connection", "user_id", userID, "error", err)
			}
		}
	}
}

// Run starts the hub's main event loop in a background goroutine
func (cs *chatsocket) Run() {
	for {
		select {
		case <-cs.hub.Ctx.Done():
			// Shutdown signal received; close all clients
			cs.hub.Mu.Lock()
			for _, client := range cs.hub.Clients {
				close(client.Send)
			}
			cs.hub.Clients = make(map[string]*entities.Client)
			cs.hub.Mu.Unlock()
			return

		case client := <-cs.hub.Register:
			cs.hub.Mu.Lock()
			// If an old connection exists for this user, we close it to prevent duplicates
			if oldClient, exists := cs.hub.Clients[client.ID]; exists {
				oldClient.Conn.Close()
				delete(cs.hub.Clients, client.ID)
				metrics.ConnectionClosed()
			}
			cs.hub.Clients[client.ID] = client
			cs.hub.Mu.Unlock()
			metrics.ConnectionOpened()
			slog.Info("websocket client registered", "user_id", client.ID, "total_connected", len(cs.hub.Clients))

		case client := <-cs.hub.UnRegister:
			cs.hub.Mu.Lock()
			if registeredClient, ok := cs.hub.Clients[client.ID]; ok && registeredClient == client {
				delete(cs.hub.Clients, client.ID)
				close(client.Send)
				metrics.ConnectionClosed()
				slog.Info("websocket client unregistered", "user_id", client.ID, "total_connected", len(cs.hub.Clients))
			}
			cs.hub.Mu.Unlock()

		case dm := <-cs.hub.DirectMessage:
			routeStart := time.Now()
			metrics.IncMessagesRouted()
			if !dm.EnqueuedAt.IsZero() {
				metrics.Observe("websocket.hub_queue_wait", time.Since(dm.EnqueuedAt))
			}
			cs.hub.Mu.RLock()
			// 1. Deliver to the recipient if they are online
			recipient, recipientOnline := cs.hub.Clients[dm.RecipientID]
			if recipientOnline && recipient != nil {
				select {
				case recipient.Send <- dm.Payload:
					// Queued successfully
				default:
					metrics.IncBufferFull()
					// Recipient buffer full, disconnect slow client
					slog.Warn("websocket send buffer full; disconnecting client", "user_id", dm.RecipientID)
					go func(c *entities.Client) {
						cs.hub.UnRegister <- c
						c.Conn.Close()
					}(recipient)
				}
			}

			// 2. Echo back to the sender's other connections (if any) to keep UI in sync
			// Only echo if the sender is different from the recipient
			if dm.SenderID != dm.RecipientID {
				sender, senderOnline := cs.hub.Clients[dm.SenderID]
				if senderOnline && sender != nil {
					select {
					case sender.Send <- dm.Payload:
						// Queued successfully
					default:
						metrics.IncBufferFull()
						slog.Warn("websocket sender buffer full", "user_id", dm.SenderID)
					}
				}
			}
			cs.hub.Mu.RUnlock()
			metrics.Observe("websocket.hub_route", time.Since(routeStart))
		}
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func readPump(c *entities.Client) {
	defer func() {
		select {
		case c.Hub.UnRegister <- c:
		case <-c.Hub.Ctx.Done():
		}
		c.Conn.Close()
	}()

	// Apply message validation / health constraints
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// SetPongHandler is called when a pong message is received from the client
	// It resets the read deadline to keep the connection alive
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure) {
				slog.Warn("websocket read failed", "user_id", c.ID, "error", err)
			}
			break
		}

	}
}

// writePump pumps messages from the hub to the WebSocket connection
func writePump(c *entities.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel, send close frame
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// Get a writer for a text message
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Optimize throughput: batch any queued messages in c.send into the same write
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			// Send periodic Ping message to the client to verify connection liveness
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (cs *chatsocket) PublishMessage(message *entities.Message) error {
	payload, err := json.Marshal(entities.WebSocketMessage{
		Message:    message.Message,
		FromUserID: message.FromUserID,
		ToUserID:   message.ToUserID,
	})
	if err != nil {
		return err
	}

	cs.hub.DirectMessage <- entities.DirectMessage{
		RecipientID: fmt.Sprintf("%s", message.ToUserID),
		SenderID:    fmt.Sprintf("%s", message.FromUserID),
		Payload:     payload,
		EnqueuedAt:  time.Now(),
	}

	return nil
}

func (cs *chatsocket) RouteToLocalConnections(event *packet.MessageEvent) {
	select {
	case cs.hub.DirectMessage <- entities.DirectMessage{
		RecipientID: event.RecipientID,
		SenderID:    event.SenderID,
		Payload:     event.Payload,
		EnqueuedAt:  event.CreatedAt,
	}:
	default:
		metrics.IncBufferFull()
		slog.Warn("local websocket hub queue full", "event_id", event.EventID)
	}
}
