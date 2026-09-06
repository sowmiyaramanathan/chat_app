package server

import (
	"backend/db/postgres"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"
)

func pingPostgres(errChan chan<- error) {
	err := postgres.Ping()
	if err != nil {
		errChan <- errors.New("postgres ping failed")
	}
}

var checks = []func(chan<- error){
	pingPostgres,
}

func Ready(w http.ResponseWriter, r *http.Request) {
	errChan := make(chan error)
	for _, f := range checks {
		go f(errChan)
	}
	for i := 0; i < len(checks); i++ {
		if err := <-errChan; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
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
