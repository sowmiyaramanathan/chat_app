package controllers

import (
	"backend/auth"
	"backend/services"
	"net/http"
)

type controller struct {
	service services.Service

	authLimiter    *auth.RateLimiter
	messageLimiter *auth.RateLimiter
}

type Controller interface {
	//user
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	RefreshToken(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
	GetNonFriends(w http.ResponseWriter, r *http.Request)

	//message
	CreateMessage(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)

	//friends
	GetMyFriends(w http.ResponseWriter, r *http.Request)
	IsFriend(w http.ResponseWriter, r *http.Request)
	SendFriendRequest(w http.ResponseWriter, r *http.Request)
	GetFriendRequests(w http.ResponseWriter, r *http.Request)
	AcceptFriendRequest(w http.ResponseWriter, r *http.Request)
	RejectFriendRequest(w http.ResponseWriter, r *http.Request)
	IsFriendRequestSent(w http.ResponseWriter, r *http.Request)
	IsRequestReceived(w http.ResponseWriter, r *http.Request)

	// websocket handler
	HandleConnection(w http.ResponseWriter, r *http.Request)
}

func New(s services.Service, authLimiter, messageLimiter *auth.RateLimiter) Controller {
	return &controller{
		service:        s,
		authLimiter:    authLimiter,
		messageLimiter: messageLimiter,
	}
}
