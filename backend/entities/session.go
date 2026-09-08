package entities

import (
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	Username     string    `json:"userName"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
