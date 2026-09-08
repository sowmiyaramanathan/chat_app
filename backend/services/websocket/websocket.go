package websocket

import (
	"backend/entities"
	"backend/entities/packet"
	redisstore "backend/redis"
	"context"

	"github.com/gorilla/websocket"
)

type chatsocket struct {
	hub      *entities.Hub
	registry redisstore.ConnectionRegistry
	instance string
}

type ChatSocket interface {
	Run()
	RunWebsocket(conn *websocket.Conn, connUserID string)
	PublishMessage(message *entities.Message) error
	RouteToLocalConnections(event *packet.MessageEvent)
}

func New(hub *entities.Hub, registry redisstore.ConnectionRegistry, instance string) ChatSocket {
	cs := &chatsocket{hub: hub, registry: registry, instance: instance}
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
