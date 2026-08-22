package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	clients = map[string]*websocket.Conn{}
	mu      = sync.Mutex{}
)

// Message structure expected from the client
type WebSocketMessage struct {
	Message    string `json:"Message"`
	FromUserID int    `json:"FromUserID"`
	ToUserID   int    `json:"ToUserID"`
}

func HandleConnection(w http.ResponseWriter, r *http.Request) {
	recipientID := chi.URLParam(r, "userID")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading socket: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	mu.Lock()
	clients[recipientID] = conn
	mu.Unlock()

	defer func() {
		mu.Lock()
		if clients[recipientID] == conn {
			delete(clients, recipientID)
		}
		mu.Unlock()
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg WebSocketMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			log.Printf("error decoding websocket message: %v", err)
			continue
		}

		// Send message to recipient if online
		recipientKey := fmt.Sprintf("%d", wsMsg.ToUserID)
		mu.Lock()
		recipientConn, recipientOnline := clients[recipientKey]
		mu.Unlock()
		if recipientOnline && recipientConn != nil {
			if err := recipientConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("write to recipient failed: %v", err)
			}
		}

		// Also send message to sender if they have a connection (keep sender's chat screen in sync)
		senderKey := fmt.Sprintf("%d", wsMsg.FromUserID)
		if senderKey != recipientKey {
			mu.Lock()
			senderConn, senderOnline := clients[senderKey]
			mu.Unlock()
			if senderOnline && senderConn != nil && senderConn != conn {
				if err := senderConn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Printf("write to sender failed: %v", err)
				}
			}
		}

		// If they are different, we echo back so the sender's UI stays in sync.
		if recipientKey != recipientID {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("echo back to self failed: %v", err)
			}
		}
	}
}
