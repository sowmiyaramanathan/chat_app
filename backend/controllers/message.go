package controllers

import (
	"backend/apperrors"
	"backend/auth"
	e "backend/entities"
	"backend/entities/packet"
	"backend/utils"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (c *controller) CreateMessage(w http.ResponseWriter, r *http.Request) {
	toID := r.URL.Query().Get("toID")
	if toID == "" {
		utils.WriteError(w, http.StatusBadRequest, "empty query params userID")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
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

	key := "user:" + claims.ID
	if !c.messageLimiter.Allow(key) {
		utils.WriteError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	}

	message.FromUserID = claims.ID
	message.ToUserID = toID

	err = c.service.Chat.CreateMessage(&message)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotAFriend) {
			utils.WriteError(w, http.StatusForbidden, "not friends")
			return
		}
		if errors.Is(err, apperrors.ErrInvalidMessage) {
			utils.WriteError(w, http.StatusUnprocessableEntity, "invalid message")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go func(message *e.Message) {
		// The request context is canceled when this handler returns. Publishing is
		// intentionally best-effort after persistence, so it needs its own deadline.
		publishCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		pubMsgEvent := packet.MessageEvent{
			EventID:     uuid.NewString(),
			MessageID:   strconv.FormatUint(uint64(message.ID), 10),
			SenderID:    message.FromUserID,
			RecipientID: message.ToUserID,
			CreatedAt:   time.Now().UTC(),
		}
		payload, marshalErr := json.Marshal(e.WebSocketMessage{
			Message:    message.Message,
			FromUserID: message.FromUserID,
			ToUserID:   message.ToUserID,
		})
		if marshalErr != nil {
			slog.Warn("failed to marshal publish message event", "error", marshalErr)
			return
		}
		pubMsgEvent.Payload = payload

		if publishErr := c.service.Redis.PublishMessage(publishCtx, pubMsgEvent); publishErr != nil {
			slog.Warn("failed to publish message event to redis", "event_id", pubMsgEvent.EventID, "error", publishErr)
		}
	}(&message)

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

	if limit < 1 {
		utils.WriteError(w, http.StatusBadRequest, "invalid query params")
		return
	} else if limit > 50 {
		limit = 50
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

	messages, err := c.service.Chat.GetMyMessages(claims.ID, toID, limit, cursor)
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
