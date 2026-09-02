package models

import (
	e "backend/entities"

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
	GetUsers(username string) ([]*p.Users, error)
	GetPublicKey(userID uint64) (string, error)

	//message
	SaveMessage(message *e.Message) (*e.Message, error)
	GetMyMessagesByFromToId(fromId, toId uint64, limit int, cursor *p.Cursor) ([]*p.Messages, error)

	//friends
	IsFriend(userAID, userBID uint64) (bool, error)
	CreateRequest(userAID, userBID uint64) error
	GetMyRequests(userID uint64) ([]*p.Requests, error)
	AcceptRequest(userAID, userBID uint64) error
	RejectRequest(userAID, userBID uint64) error
	IsRequestSent(userAID, userBID uint64) (bool, error)
	IsRequestReceived(userAID, userBID uint64) (bool, error)
}

func New(Db *gorm.DB) Model {
	return &model{
		Db: Db,
	}
}
