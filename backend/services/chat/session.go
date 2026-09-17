package chat

import (
	e "backend/entities"
	"time"
)

func (c *chat) CreateSession(session *e.Session) error {
	err := c.m.CreateSession(session)
	if err != nil {
		return err
	}

	return nil
}

func (c *chat) RotateSessionRefreshToken(ID, currentToken, nextToken string, expiresAt time.Time) (bool, error) {
	return c.m.RotateSessionRefreshToken(ID, currentToken, nextToken, expiresAt)
}

func (c *chat) GetSessionByRefreshToken(token string) (*e.Session, error) {
	session, err := c.m.GetSessionByRefreshToken(token)
	if err != nil {
		return nil, err
	}

	return session, nil
}
