package services

import (
	"backend/models"
	"backend/services/chat"
	cs "backend/services/websocket"
)

type Service struct {
	Chat chat.Chat
	CS   cs.ChatSocket
}

func New(m models.Model, cs cs.ChatSocket) Service {
	s := Service{
		Chat: chat.New(&m),
		CS:   cs,
	}

	return s
}
