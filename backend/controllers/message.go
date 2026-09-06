package controllers

import (
	"backend/apperrors"
	"backend/auth"
	e "backend/entities"
	"backend/entities/packet"
	"backend/metrics"
	"backend/utils"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (c *controller) CreateMessage(w http.ResponseWriter, r *http.Request) {
	if !c.bucket.Take(1) {
		utils.WriteError(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	toID := r.URL.Query().Get("toID")
	if toID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	var message e.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, "invalid json")
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

	message.FromUserID = claims.ID
	message.ToUserID = toID

	err = c.s.Chat.CreateMessage(&message)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	publishStart := time.Now()
	if err := c.s.CS.PublishMessage(&message); err != nil {
		metrics.Observe("websocket.publish_enqueue", time.Since(publishStart))
		slog.Error("could not publish websocket message", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}
	metrics.Observe("websocket.publish_enqueue", time.Since(publishStart))

	w.Write([]byte("Message Sent"))
}

func (c *controller) GetMessages(w http.ResponseWriter, r *http.Request) {
	toID := r.URL.Query().Get("toID")
	if toID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid query params")
		return
	}

	cursorStr := r.URL.Query().Get("cursor")
	cursor, err := utils.DecodeCursor(cursorStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid query params")
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

	messages, err := c.s.Chat.GetMyMessages(claims.ID, toID, limit, cursor)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "internal error")
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
