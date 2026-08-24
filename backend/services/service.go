package services

import (
	e "backend/entities"
	"backend/models"
	cs "backend/services/websocket"

	"github.com/gorilla/websocket"
)

type service struct {
	m  models.Model
	ws cs.ChatSocket
}

type Service interface {
	//user
	CreateUser(user *e.User) error
	LoginUser(username, password string) (uint, error)
	GetAllUsers(username string) (*[]e.Users, error)
	GetPublicKey(userID uint64) (string, error)

	//message
	CreateMessage(message *e.Message) error
	GetMyMessages(fromId, toId uint64) (*[]e.Messages, error)

	//friends
	CheckIsFriend(userAID, userBID uint64) (bool, error)
	CreateFriendRequest(userAID, userBID uint64) error
	GetFriendRequests(userID uint64) ([]e.Requests, error)
	AcceptFriendRequest(userAID, userBID uint64) error
	RejecttFriendRequest(userAID, userBID uint64) error
	CheckIsFriendRequestSent(userAID, userBID uint64) (bool, error)
	CheckIsFriendRequestReceived(userAID, UserBID uint64) (bool, error)

	// websocket
	RunWebsocket(conn *websocket.Conn, connUserID string)
}

func New(m models.Model, ws cs.ChatSocket) Service {
	return &service{
		m:  m,
		ws: ws,
	}
}
