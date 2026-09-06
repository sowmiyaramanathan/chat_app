package chat

import (
	e "backend/entities"
	p "backend/entities/packet"
	"backend/metrics"
	"strings"
	"time"
)

func prepareMessage(message *e.Message) {
	message.Message = strings.TrimSpace(message.Message)
}

func (c *chat) CreateMessage(message *e.Message) error {
	prepareMessage(message)
	dbStart := time.Now()
	_, err := c.m.SaveMessage(message)
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
