package models

import (
	e "backend/entities"
	"log/slog"
)

func (m *model) CreateSession(session *e.Session) error {
	err := m.Db.Create(session).Error
	if err != nil {
		slog.Debug("error creating session", "error", err)
		return err
	}

	return nil
}

func (m *model) UpdateSessionRefreshToken(ID, token string) error {
	err := m.Db.Model(&e.Session{}).Where("id = ?", ID).Update("refresh_token", token).Error
	if err != nil {
		slog.Debug("error updating refresh token", "error", err)
		return err
	}

	return nil
}

func (m *model) GetSessionByRefreshToken(token string) (*e.Session, error) {
	session := &e.Session{}

	err := m.Db.Model(&e.Session{}).Where("refresh_token = ?", token).Find(session).Error
	if err != nil {
		slog.Debug("error getting session by refresh token", "error", err)
		return nil, err
	}

	return session, nil
}
