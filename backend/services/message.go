package services

import (
	e "backend/entities"
	"errors"
	"strings"
)

func prepareMessage(message *e.Message) {
	message.Message = strings.TrimSpace(message.Message)
}

func (s *service) CreateMessage(message *e.Message) error {
	prepareMessage(message)
	_, err := s.m.SaveMessage(message)
	if err != nil {
		return errors.New("could not create message")
	}

	return nil
}

func (s *service) GetMyMessages(fromId, toId uint64) (*[]e.Messages, error) {
	messages, err := s.m.GetMyMessagesByFromToId(fromId, toId)

	if err != nil {
		return nil, errors.New("could not get messages")
	}

	return messages, nil
}
