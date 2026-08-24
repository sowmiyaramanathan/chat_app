package controllers

import (
	"backend/services"
	"net/http"
)

type controller struct {
	s services.Service
}

type Controller interface {
	//user
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
	GetAllUsers(w http.ResponseWriter, r *http.Request)
	GetPublicKey(w http.ResponseWriter, r *http.Request)

	//message
	CreateMessage(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)

	//friends
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

func New(s services.Service) Controller {
	return &controller{
		s: s,
	}
}
