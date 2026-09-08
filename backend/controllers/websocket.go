package controllers

import (
	"backend/auth"
	"backend/utils"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range websocketOrigins() {
			if origin == allowedOrigin {
				return true
			}
		}
		return false
	},
}

func websocketOrigins() []string {
	configured := os.Getenv("FRONTEND_ORIGINS")
	if configured == "" {
		configured = os.Getenv("ALLOWED_ORIGIN")
	}
	if configured == "" {
		return []string{"http://localhost:3000", "http://localhost:3001"}
	}

	origins := strings.Split(configured, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}

func (c *controller) HandleConnection(w http.ResponseWriter, r *http.Request) {
	connectedUserID := chi.URLParam(r, "userID")
	if connectedUserID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty URL params user")
		return
	}

	// 1. AUTHENTICATION & JWT CLAIM VALIDATION
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params token")
		return
	}

	token, err := jwtauth.VerifyRequest(auth.TokenAuth, r, func(r *http.Request) string {
		return tokenString
	})
	if err != nil || token == nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	rawUserID, ok := token.Get("userID")
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "missing userID claim")
		return
	}

	claimUserID := fmt.Sprintf("%v", rawUserID)
	if claimUserID != connectedUserID {
		utils.WriteError(w, http.StatusForbidden, "forbidden")
		return
	}

	// 2. UPGRADE TO WEBSOCKET
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		utils.WriteError(w, http.StatusBadRequest, "WebSocket upgrade failed")
		return
	}

	c.s.CS.RunWebsocket(conn, connectedUserID)
}
