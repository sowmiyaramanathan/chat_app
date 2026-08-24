package websocket

import (
	"backend/entities"
	"context"

	"github.com/gorilla/websocket"
)

type chatsocket struct {
	hub *entities.Hub
}

type ChatSocket interface {
	Run()
	RunWebsocket(conn *websocket.Conn, connUserID string)
}

func New(hub *entities.Hub) ChatSocket {
	cs := &chatsocket{hub: hub}
	go cs.Run()
	return cs
}

// NewHub creates and initializes a new Hub
func NewHub() *entities.Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &entities.Hub{
		Clients:       make(map[string]*entities.Client),
		DirectMessage: make(chan entities.DirectMessage, 256),
		Register:      make(chan *entities.Client),
		UnRegister:    make(chan *entities.Client),
		Ctx:           ctx,
		Cancel:        cancel,
	}
}
