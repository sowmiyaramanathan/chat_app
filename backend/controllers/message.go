package controllers

import (
	"backend/auth"
	e "backend/entities"
	"backend/entities/packet"
	"backend/utils"
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

	w.Write([]byte("Message Sent"))
}

func (c *controller) GetMessages(w http.ResponseWriter, r *http.Request) {
	toIdStr := r.URL.Query().Get("to_id")
	toId, err := strconv.ParseUint(toIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	cursorStr := r.URL.Query().Get("cursor")
	cursor, err := utils.DecodeCursor(cursorStr)
	if err != nil {
		http.Error(w, "Invalid cursor", http.StatusBadRequest)
		return
	}

	claims := auth.ExtractToken(r)

	messages, err := c.s.GetMyMessages(uint64(claims.Id), toId, limit, cursor)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hasNextPage := len(messages) > limit
	if hasNextPage {
		messages = messages[:limit]
	}

	response := packet.AllMessagesRes{
		Messages: messages,
		PageInfo: packet.PageInfo{HasNextPage: hasNextPage},
	}
	if response.PageInfo.HasNextPage {
		lastMsg := messages[len(messages)-1]
		response.PageInfo.EndCursor = utils.EncodeCursor(lastMsg.CreatedAt, lastMsg.ID)
	}

	json.NewEncoder(w).Encode(response)
}
