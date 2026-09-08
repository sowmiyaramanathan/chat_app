package controllers

import (
	"backend/apperrors"
	"backend/auth"
	e "backend/entities"
	"backend/utils"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

func (c *controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if !c.bucket.Take(1) {
		slog.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
		utils.WriteError(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	var user e.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "invalid json")
		return
	}

	err = c.service.Chat.CreateUser(&user)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			utils.WriteError(w, http.StatusConflict, "username")
			return
		} else if errors.Is(err, apperrors.ErrMobileAlreadyExists) {
			utils.WriteError(w, http.StatusConflict, "number")
			return
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "internal error")
			return
		}
	}

	w.Write([]byte("User Created Successfully"))
}

func (c *controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	if !c.bucket.Take(1) {
		slog.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
		utils.WriteError(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	var user e.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	userID, err := c.service.Chat.LoginUser(user.UserName, user.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			utils.WriteError(w, http.StatusConflict, "username")
			return
		}
		if errors.Is(err, apperrors.ErrInvalidCredentials) {
			utils.WriteError(w, http.StatusConflict, "password")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	// Set access/refresh token expiry durations
	accessTokenDuration := time.Minute * 15
	refreshTokenDuration := time.Hour * 24 * 7 // 7 days

	accessToken, refreshToken, err := auth.CreateToken(userID, user.UserName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	now := time.Now()
	session := &e.Session{
		ID:           userID,
		Username:     user.UserName,
		RefreshToken: refreshToken,
		ExpiresAt:    now.Add(refreshTokenDuration),
	}

	if err := c.service.Chat.CreateSession(session); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":                "Logged in successfully",
		"token":                 accessToken,
		"refreshToken":          refreshToken,
		"accessTokenExpiresIn":  int(accessTokenDuration.Seconds()),
		"refreshTokenExpiresIn": int(refreshTokenDuration.Seconds()),
	})
}

func (c *controller) Profile(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"name": claims.Username,
	})
}

func (c *controller) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	users, err := c.service.Chat.GetAllUsers(claims.Username)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	data := users
	json.NewEncoder(w).Encode(data)
}

func (c *controller) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Validate the refresh token exists and is not expired
	session, err := c.service.Chat.GetSessionByRefreshToken(req.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	if time.Now().After(session.ExpiresAt) {
		utils.WriteError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	// Query the user by username from the session
	user, err := c.service.Chat.GetUserByUsername(session.Username)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "User not found for refresh token")
		return
	}

	// Create new tokens
	accessToken, _, err := auth.CreateToken(user.ID, user.UserName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create tokens")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":       "Token refreshed successfully",
		"token":        accessToken,
		"refreshToken": req.RefreshToken,
	})
}
