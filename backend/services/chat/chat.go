package chat

import (
	e "backend/entities"
	p "backend/entities/packet"
	"backend/models"
	"time"
)

type chat struct {
	m models.Model
}

type Chat interface {
	//user
	CreateUser(user *e.User) error
	LoginUser(username, password string) (string, error)
	GetNonFriends(userID string, limit int, cursor string) ([]*p.Users, error)
	GetUserByUsername(username string) (*e.User, error)

	//message
	CreateMessage(message *e.Message) error
	GetMyMessages(fromId, toId string, limit int, cursor *p.Cursor) ([]*p.Messages, error)

	//friends
	GetMyFriends(userID string, limit int, cursor string) ([]*p.Users, error)
	CheckIsFriend(userAID, userBID string) (bool, error)
	CreateFriendRequest(userAID, userBID string) error
	GetFriendRequests(userID string) ([]*p.Requests, error)
	AcceptFriendRequest(userAID, userBID string) error
	RejecttFriendRequest(userAID, userBID string) error
	CheckIsFriendRequestSent(userAID, userBID string) (bool, error)
	CheckIsFriendRequestReceived(userAID, UserBID string) (bool, error)

	// session
	CreateSession(session *e.Session) error
	RotateSessionRefreshToken(ID, currentToken, nextToken string, expiresAt time.Time) (bool, error)
	GetSessionByRefreshToken(token string) (*e.Session, error)
}

func New(m *models.Model) Chat {
	return &chat{
		m: *m,
	}
}
