package main

import (
	"backend/server"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("failed to load environment", "error", err)
		os.Exit(1)
	} else {
		slog.Info("environment loaded")
	}

	server.Run()
}
