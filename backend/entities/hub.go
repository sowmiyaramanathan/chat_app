package entities

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketMessage represents the message structure expected from/to the client
type WebSocketMessage struct {
	Message    string `json:"Message"`
	FromUserID string `json:"FromUserID"`
	ToUserID   string `json:"ToUserID"`
}

// DirectMessage represents a routed direct 1-to-1 message in the Hub
type DirectMessage struct {
	RecipientID string
	SenderID    string
	Payload     []byte
	EnqueuedAt  time.Time
}

// Client represents a single connected user's WebSocket connection
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte // Buffered channel for outbound messages
	ID   string      // The authenticated User ID
}

// Hub maintains the set of active clients and handles 1-to-1 message routing
type Hub struct {
	// Registered clients mapped by their User ID
	Clients map[string]*Client

	// Channel for incoming 1-to-1 direct messages
	DirectMessage chan DirectMessage

	// Channel for client registration
	Register chan *Client

	// Channel for client unregistration
	UnRegister chan *Client

	// Context for graceful shutdown
	Ctx    context.Context
	Cancel context.CancelFunc
	Mu     sync.RWMutex
}
