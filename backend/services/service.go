package services

import (
	"backend/models"
	"backend/redis"
	"backend/services/chat"
	cs "backend/services/websocket"
)

type Service struct {
	Chat chat.Chat
	CS   cs.ChatSocket

	Redis redis.RedisPubSub
}

func New(m models.Model, cs cs.ChatSocket, redis redis.RedisPubSub) Service {
	s := Service{
		Chat:  chat.New(&m),
		CS:    cs,
		Redis: redis,
	}

	return s
}
