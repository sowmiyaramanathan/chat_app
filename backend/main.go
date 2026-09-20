package main

import (
	"backend/auth"
	"backend/server"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := godotenv.Load(); err == nil {
		slog.Info("environment loaded")
	}
	if err := auth.Init(); err != nil {
		slog.Error("authentication configuration failed", "error", err)
		os.Exit(1)
	}

	server.Run()
}
