package models

import (
	e "backend/entities"
	"log/slog"
	"time"
)

func (m *model) CreateSession(session *e.Session) error {
	err := m.Db.Create(session).Error
	if err != nil {
		slog.Debug("error creating session", "error", err)
		return err
	}

	return nil
}

func (m *model) RotateSessionRefreshToken(ID, currentToken, nextToken string, expiresAt time.Time) (bool, error) {
	result := m.Db.Model(&e.Session{}).
		Where("id = ? AND refresh_token = ?", ID, currentToken).
		Updates(map[string]interface{}{"refresh_token": nextToken, "expires_at": expiresAt})
	if result.Error != nil {
		slog.Debug("error rotating refresh token", "error", result.Error)
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

func (m *model) GetSessionByRefreshToken(token string) (*e.Session, error) {
	session := &e.Session{}

	err := m.Db.Where("refresh_token = ?", token).Take(session).Error
	if err != nil {
		slog.Debug("error getting session by refresh token", "error", err)
		return nil, err
	}

	return session, nil
}
