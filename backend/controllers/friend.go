package controllers

import (
	"backend/apperrors"
	"backend/auth"
	"backend/entities/packet"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"
)

func (c *controller) IsFriend(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	resp, err := c.service.Chat.CheckIsFriend(claims.ID, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if !resp {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": true})
}

func (c *controller) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	err = c.service.Chat.CreateFriendRequest(claims.ID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidFriendRequest) {
			utils.WriteError(w, http.StatusBadRequest, "invalid friend request")
			return
		}
		if errors.Is(err, apperrors.ErrFriendRequestExists) {
			utils.WriteError(w, http.StatusConflict, "friend request already exists")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": "Sent"})
}

func (c *controller) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	requests, err := c.service.Chat.GetFriendRequests(claims.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	data := requests
	json.NewEncoder(w).Encode(data)
}

func (c *controller) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	err = c.service.Chat.AcceptFriendRequest(userID, claims.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Write([]byte("Friend Request Accepted"))
}

func (c *controller) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	err = c.service.Chat.RejecttFriendRequest(userID, claims.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.Write([]byte("Friend Request Rejected"))
}

func (c *controller) IsFriendRequestSent(w http.ResponseWriter, r *http.Request) {
	toUserID := r.URL.Query().Get("userID")

	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	resp, err := c.service.Chat.CheckIsFriendRequestSent(claims.ID, toUserID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if !resp {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": true})
}

func (c *controller) IsRequestReceived(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	resp, err := c.service.Chat.CheckIsFriendRequestReceived(userID, claims.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if !resp {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": true})
}

func (c *controller) GetMyFriends(w http.ResponseWriter, r *http.Request) {
	claims, err := auth.ExtractToken(r)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidToken) {
			utils.WriteError(w, http.StatusUnauthorized, "Unauthorized - claims missing")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "Failed to extract token")
		return
	}

	limit, cursor, err := userPageParams(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid query params")
		return
	}

	users, err := c.service.Chat.GetMyFriends(claims.ID, limit, cursor)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	json.NewEncoder(w).Encode(userPage(users, limit))
}

func userPage(users []*packet.Users, limit int) packet.UsersPage {
	page := packet.UsersPage{Users: users, PageInfo: packet.PageInfo{HasNextPage: len(users) > limit}}
	if page.PageInfo.HasNextPage {
		page.Users = users[:limit]
	}
	if len(page.Users) > 0 {
		page.PageInfo.EndCursor = page.Users[len(page.Users)-1].ID
	}
	return page
}
