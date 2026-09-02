package services

import (
	e "backend/entities"
	"backend/metrics"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func prepareMessage(message *e.Message) {
	message.Message = strings.TrimSpace(message.Message)
}

func (s *service) CreateMessage(message *e.Message) error {
	prepareMessage(message)
	dbStart := time.Now()
	_, err := s.m.SaveMessage(message)
	metrics.Observe("db.save_message", time.Since(dbStart))
	if err != nil {
		return errors.New("could not create message")
	}

	publishStart := time.Now()
	if err := s.cs.PublishMessage(message); err != nil {
		metrics.Observe("websocket.publish_enqueue", time.Since(publishStart))
		return fmt.Errorf("could not publish websocket message: %w", err)
	}
	metrics.Observe("websocket.publish_enqueue", time.Since(publishStart))

	return nil
}

func (s *service) GetMyMessages(fromId, toId uint64) (*[]e.Messages, error) {
	messages, err := s.m.GetMyMessagesByFromToId(fromId, toId)

	if err != nil {
		return nil, errors.New("could not get messages")
	}

	return messages, nil
}

func (s *service) RunWebsocket(conn *websocket.Conn, connUserID string) {
	s.cs.RunWebsocket(conn, connUserID)
}

func (s *service) PublishMessage(message *e.Message) error {
	return s.cs.PublishMessage(message)
}
