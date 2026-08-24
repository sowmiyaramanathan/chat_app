package controllers

import (
	"backend/auth"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *controller) HandleConnection(w http.ResponseWriter, r *http.Request) {
	// 1. AUTHENTICATION & JWT CLAIM VALIDATION
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	token, err := jwtauth.VerifyRequest(auth.TokenAuth, r, func(r *http.Request) string {
		return r.URL.Query().Get("token")
	})
	if err != nil || token == nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	rawUserID, ok := token.Get("userID")
	if !ok {
		http.Error(w, "missing userID claim", http.StatusUnauthorized)
		return
	}

	claimUserID := fmt.Sprintf("%v", rawUserID)
	connectedUserID := chi.URLParam(r, "userID")
	if connectedUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if claimUserID != connectedUserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 2. UPGRADE TO WEBSOCKET
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading socket: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	c.s.RunWebsocket(conn, connectedUserID)
}
