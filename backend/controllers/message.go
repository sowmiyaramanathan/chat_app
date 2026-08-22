package controllers

import (
	"backend/auth"
	e "backend/entities"
	"encoding/json"
	"net/http"
	"strconv"
)

func (c *controller) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message e.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Cannot Process your Request", http.StatusUnprocessableEntity)
		return
	}
	toIdString := r.URL.Query().Get("to_id")
	toId, err := strconv.ParseUint(toIdString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)
	message.FromUserID = uint64(claims.Id)
	message.ToUserID = toId

	err = c.s.CreateMessage(&message)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Message Created"))
}

func (c *controller) GetMessages(w http.ResponseWriter, r *http.Request) {
	toIdString := r.URL.Query().Get("to_id")
	toId, err := strconv.ParseUint(toIdString, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)

	messages, err := c.s.GetMyMessages(uint64(claims.Id), toId)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := messages
	json.NewEncoder(w).Encode(data)
}
