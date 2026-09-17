package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"not null;index" json:"userId"`
	Username     string    `json:"userName"`
	RefreshToken string    `gorm:"not null;uniqueIndex" json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewSession(userID, username, refreshToken string, expiresAt time.Time) *Session {
	return &Session{ID: uuid.NewString(), UserID: userID, Username: username, RefreshToken: refreshToken, ExpiresAt: expiresAt}
}
