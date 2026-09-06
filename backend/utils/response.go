package utils

import (
	"backend/entities/packet"
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(packet.ErrorResponse{
		Error: code,
	})
}
