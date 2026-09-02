package models

import (
	e "backend/entities"
)

func (m *model) SaveMessage(message *e.Message) (*e.Message, error) {
	err := m.Db.Create(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (m *model) GetMyMessagesByFromToId(fromId, toId uint64) (*[]e.Messages, error) {
	messages := &[]e.Messages{}

	err := m.Db.
		Model(&e.Message{}).Order("created_at ASC").
		Select("id", "from_user_id", "to_user_id", "message", "created_at").
		Where("(from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?)", fromId, toId, toId, fromId).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}
	return messages, nil
}
