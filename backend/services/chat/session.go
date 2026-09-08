package chat

import e "backend/entities"

func (c *chat) CreateSession(session *e.Session) error {
	err := c.m.CreateSession(session)
	if err != nil {
		return err
	}

	return nil
}

func (c *chat) UpdateSessionRefreshToken(ID, token string) error {
	err := c.m.UpdateSessionRefreshToken(ID, token)
	if err != nil {
		return err
	}

	return nil
}

func (c *chat) GetSessionByRefreshToken(token string) (*e.Session, error) {
	session, err := c.m.GetSessionByRefreshToken(token)
	if err != nil {
		return nil, err
	}

	return session, nil
}
