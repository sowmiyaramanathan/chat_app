package chat

import (
	"backend/apperrors"
	e "backend/entities"
	p "backend/entities/packet"
	"html"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Debug("failed to hash password", "error", err)
		return "", err
	}
	return string(bytes), nil
}

func prepareUser(user *e.User) error {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.UserName = html.EscapeString(strings.TrimSpace(user.UserName))
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return nil
}

func (c *chat) CreateUser(user *e.User) error {
	err := prepareUser(user)
	if err != nil {
		return err
	}

	_, err = c.m.GetUserByUsername(user.UserName)
	if err == nil {
		return apperrors.ErrUserAlreadyExists
	}

	_, err = c.m.GetUserByMobilenumber(user.Mobilenumber)
	if err == nil {
		return apperrors.ErrMobileAlreadyExists
	}

	user.ID = uuid.NewString()
	_, err = c.m.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (c *chat) LoginUser(username, password string) (ID string, err error) {
	user, err := c.m.GetUserByUsername(username)
	if err != nil {
		return ID, apperrors.ErrUserNotFound
	}
	err = verifyPassword(user.Password, password)
	if err != nil {
		return ID, apperrors.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (c *chat) GetNonFriends(userID string, limit int, cursor string) ([]*p.Users, error) {
	users, err := c.m.GetNonFriends(userID, limit, cursor)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *chat) GetUserByUsername(username string) (*e.User, error) {
	user, err := c.m.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
