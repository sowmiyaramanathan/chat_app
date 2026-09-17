package models

import (
	e "backend/entities"
	"time"

	p "backend/entities/packet"

	"gorm.io/gorm"
)

type model struct {
	Db *gorm.DB
}

type Model interface {
	//user
	SaveUser(user *e.User) (*e.User, error)
	// GetUserById(id uint64) (*e.User, error)
	GetUserByUsername(username string) (*e.User, error)
	GetUserByMobilenumber(number string) (*e.User, error)
	GetNonFriends(userID string, limit int, cursor string) ([]*p.Users, error)

	//message
	SaveMessage(message *e.Message) (*e.Message, error)
	GetMyMessagesByFromToId(fromId, toId string, limit int, cursor *p.Cursor) ([]*p.Messages, error)

	//friends
	GetMyFriends(userID string, limit int, cursor string) ([]*p.Users, error)
	IsFriend(userAID, userBID string) (bool, error)
	CreateRequest(userAID, userBID string) error
	GetMyRequests(userID string) ([]*p.Requests, error)
	AcceptRequest(userAID, userBID string) error
	RejectRequest(userAID, userBID string) error
	IsRequestSent(userAID, userBID string) (bool, error)
	IsRequestReceived(userAID, userBID string) (bool, error)

	// session
	CreateSession(session *e.Session) error
	RotateSessionRefreshToken(ID, currentToken, nextToken string, expiresAt time.Time) (bool, error)
	GetSessionByRefreshToken(token string) (*e.Session, error)
}

func New(Db *gorm.DB) Model {
	return &model{
		Db: Db,
	}
}
