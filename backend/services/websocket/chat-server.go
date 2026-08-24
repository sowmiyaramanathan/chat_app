package websocket

import (
	"backend/entities"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	// Maximum message size allowed from client (e.g., 512KB)
	maxMessageSize = 512 * 1024
)

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
			}
			cs.hub.Clients[client.ID] = client
			cs.hub.Mu.Unlock()
			log.Printf("User %s registered (total connected: %d)", client.ID, len(cs.hub.Clients))

		case client := <-cs.hub.UnRegister:
			cs.hub.Mu.Lock()
			if registeredClient, ok := cs.hub.Clients[client.ID]; ok && registeredClient == client {
				delete(cs.hub.Clients, client.ID)
				close(client.Send)
				log.Printf("User %s unregistered (total connected: %d)", client.ID, len(cs.hub.Clients))
			}
			cs.hub.Mu.Unlock()

		case dm := <-cs.hub.DirectMessage:
			cs.hub.Mu.RLock()
			// 1. Deliver to the recipient if they are online
			recipient, recipientOnline := cs.hub.Clients[dm.RecipientID]
			if recipientOnline && recipient != nil {
				select {
				case recipient.Send <- dm.Payload:
					// Queued successfully
				default:
					// Recipient buffer full, disconnect slow client
					log.Printf("User %s send buffer full, disconnecting", dm.RecipientID)
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
						log.Printf("User %s (sender) buffer full", dm.SenderID)
					}
				}
			}
			cs.hub.Mu.RUnlock()
		}
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func readPump(c *entities.Client) {
	defer func() {
		c.Hub.UnRegister <- c
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
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure) {
				log.Printf("Read error for user %s: %v", c.ID, err)
			}
			break
		}

		// Decode the message to read recipient details
		var wsMsg entities.WebSocketMessage
		if err := json.Unmarshal(payload, &wsMsg); err != nil {
			log.Printf("Invalid message format from user %s: %v", c.ID, err)
			continue
		}

		// Routing verification (ensure sender claim matches payload)
		senderKey := fmt.Sprintf("%d", wsMsg.FromUserID)
		if senderKey != c.ID {
			log.Printf("Warning: User %s attempted to send message as user %d (mismatch)", c.ID, wsMsg.FromUserID)
			continue
		}

		// Route message via the hub's directMessage channel
		recipientKey := fmt.Sprintf("%d", wsMsg.ToUserID)
		c.Hub.DirectMessage <- entities.DirectMessage{
			RecipientID: recipientKey,
			SenderID:    senderKey,
			Payload:     payload,
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

	// RUN READ AND WRITE PUMPS IN SEPARATE GOROUTINES
	// This separates reading from writing, resolving concurrent-write safety issues.
	go writePump(client)
	go readPump(client)
}
