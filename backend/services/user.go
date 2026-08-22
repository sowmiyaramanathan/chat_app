package services

import (
	e "backend/entities"
	"errors"
	"fmt"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func prepareUser(user *e.User) {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Password = html.EscapeString(strings.TrimSpace(user.Password))
	user.Password = hashPassword(user.Password)
}

func (s *service) CreateUser(user *e.User) error {
	prepareUser(user)
	_, err := s.m.GetUserByUsername(user.Username)
	fmt.Println("username", err)
	if err == nil {
		return errors.New("username")
	}

	_, err = s.m.GetUserByMobilenumber(user.Mobilenumber)
	if err == nil {
		return errors.New("mobile_number")
	}

	_, err = s.m.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *service) LoginUser(username, password string) (uint, error) {
	user, err := s.m.GetUserByUsername(username)
	if err != nil {
		return 0, errors.New("username")
	}
	err = verifyPassword(user.Password, password)
	if err != nil {
		return 0, errors.New("wrong_password")
	}

	return user.ID, nil
}

func (s *service) GetAllUsers(username string) (*[]e.Users, error) {
	users, err := s.m.GetUsers(username)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) GetPublicKey(userID uint64) (string, error) {
	pubKey, err := s.m.GetPublicKey(userID)
	if err != nil {
		return "", err
	}
	return pubKey, nil
}
