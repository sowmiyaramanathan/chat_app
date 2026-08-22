package controllers

import (
	"backend/auth"
	"encoding/json"
	"net/http"
	"strconv"
)

func (c *controller) IsFriend(w http.ResponseWriter, r *http.Request) {
	withUserIdString := r.URL.Query().Get("with_user_id")
	withUser, err := strconv.ParseUint(withUserIdString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)

	resp, err := c.s.CheckIsFriend(uint64(claims.Id), withUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !resp {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": true})
}

func (c *controller) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	toUserIDString := r.URL.Query().Get("to_user_id")
	toUserID, err := strconv.ParseUint(toUserIDString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)
	err = c.s.CreateFriendRequest(uint64(claims.Id), toUserID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": "Sent"})
}

func (c *controller) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	claims := auth.ExtractToken(r)

	requests, err := c.s.GetFriendRequests(uint64(claims.Id))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := requests
	json.NewEncoder(w).Encode(data)
}

func (c *controller) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	fromUserIDString := r.URL.Query().Get("from_user_id")
	fromUserID, err := strconv.ParseUint(fromUserIDString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)
	err = c.s.AcceptFriendRequest(fromUserID, uint64(claims.Id))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Friend Request Accepted"))
}

func (c *controller) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	fromUserIDString := r.URL.Query().Get("from_user_id")
	fromUserID, err := strconv.ParseUint(fromUserIDString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)
	err = c.s.RejecttFriendRequest(fromUserID, uint64(claims.Id))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Friend Request Rejected"))
}

func (c *controller) IsFriendRequestSent(w http.ResponseWriter, r *http.Request) {
	toUserIDString := r.URL.Query().Get("to_user_id")
	toUserID, err := strconv.ParseUint(toUserIDString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)
	resp, err := c.s.CheckIsFriendRequestSent(uint64(claims.Id), toUserID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !resp {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": true})
}

func (c *controller) IsRequestReceived(w http.ResponseWriter, r *http.Request) {
	fromUserIDString := r.URL.Query().Get("from_user_id")
	fromUserID, err := strconv.ParseUint(fromUserIDString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)
	resp, err := c.s.CheckIsFriendRequestReceived(fromUserID, uint64(claims.Id))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !resp {
		json.NewEncoder(w).Encode(map[string]interface{}{"data": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"data": true})
}
