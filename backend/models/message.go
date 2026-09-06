package models

import (
	e "backend/entities"
	p "backend/entities/packet"
)

func (m *model) SaveMessage(message *e.Message) (*e.Message, error) {
	err := m.Db.Create(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *model) GetMyMessagesByFromToId(fromId, toId string, limit int, cursor *p.Cursor) (messages []*p.Messages, err error) {
	messages = []*p.Messages{}

	query := m.Db.Model(&e.Message{}).
		Select("id", "from_user_id", "to_user_id", "message", "created_at").
		Where("(from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?)", fromId, toId, toId, fromId)

	if cursor != nil {
		query = query.Where("(created_at, id) < (?, ?)", cursor.CreatedAt, cursor.ID)
	}

	err = query.Order("created_at DESC, id DESC").Limit(limit + 1).Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
