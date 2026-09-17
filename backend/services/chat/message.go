package chat

import (
	"backend/apperrors"
	e "backend/entities"
	p "backend/entities/packet"
	"backend/metrics"
	"log/slog"
	"strings"
	"time"
)

const maxMessageLength = 4000

func prepareMessage(message *e.Message) {
	message.Message = strings.TrimSpace(message.Message)
}

func (c *chat) CreateMessage(message *e.Message) error {
	prepareMessage(message)
	if message.Message == "" || len(message.Message) > maxMessageLength {
		return apperrors.ErrInvalidMessage
	}

	isFrd, err := c.CheckIsFriend(message.FromUserID, message.ToUserID)
	if err != nil {
		return err
	}
	if !isFrd {
		slog.Info("attempting to send message to an unknown person")
		return apperrors.ErrNotAFriend
	}

	dbStart := time.Now()
	_, err = c.m.SaveMessage(message)
	metrics.Observe("db.save_message", time.Since(dbStart))
	if err != nil {
		return err
	}

	return nil
}

func (c *chat) GetMyMessages(fromId, toId string, limit int, cursor *p.Cursor) ([]*p.Messages, error) {
	messages, err := c.m.GetMyMessagesByFromToId(fromId, toId, limit, cursor)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
