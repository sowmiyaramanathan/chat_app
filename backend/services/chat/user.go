package chat

import (
	"backend/apperrors"
	e "backend/entities"
	p "backend/entities/packet"
	"html"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func prepareUser(user *e.User) {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.UserName = html.EscapeString(strings.TrimSpace(user.UserName))
	user.Password = html.EscapeString(strings.TrimSpace(user.Password))
	user.Password = hashPassword(user.Password)
}

func (c *chat) CreateUser(user *e.User) error {
	prepareUser(user)
	_, err := c.m.GetUserByUsername(user.UserName)
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

func (c *chat) GetAllUsers(username string) ([]*p.Users, error) {
	users, err := c.m.GetUsers(username)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (c *chat) GetPublicKey(userID string) (string, error) {
	pubKey, err := c.m.GetPublicKey(userID)
	if err != nil {
		return "", err
	}
	return pubKey, nil
}
