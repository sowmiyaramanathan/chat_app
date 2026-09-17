package server

import (
	"backend/db/postgres"
	"backend/redis"
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

func Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := postgres.PingContext(ctx); err != nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "postgres unavailable")
		return
	}
	if err := redis.PingContext(ctx); err != nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "redis unavailable")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{"message": "OK"})
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "UP",
	})
}
