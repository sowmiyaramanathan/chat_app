package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"backend/auth"
	"backend/controllers"
	"backend/db/postgres"
	"backend/entities/packet"
	"backend/metrics"
	"backend/models"
	"backend/redis"
	"backend/routes"
	"backend/services"
	"backend/services/websocket"

	"github.com/google/uuid"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	Db := postgres.ConnectDatabase()
	sqlDB, err := Db.DB()
	if err != nil {
		slog.Error("failed to get database connection pool", "error", err)
		return
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats := sqlDB.Stats()
				slog.Info("database pool metrics", "open", stats.OpenConnections, "in_use", stats.InUse, "idle", stats.Idle, "wait_count", stats.WaitCount, "wait_duration", stats.WaitDuration)
			case <-ctx.Done():
				return
			}
		}
	}()

	hub := websocket.NewHub()

	redisClient, err := redis.New(ctx)
	if err != nil {
		slog.Error("failed to start redis", "error", err)
		return
	}
	defer redisClient.Close()
	instanceID := envString("INSTANCE_ID", uuid.NewString())
	ws := websocket.New(hub, redisClient, instanceID)

	go func() {
		err := redisClient.SubscribeMessages(ctx, func(event packet.MessageEvent) {
			slog.Info("redis message received", "instance_id", instanceID, "event_id", event.EventID, "recipient_id", event.RecipientID)
			ws.RouteToLocalConnections(&event)
		})
		if err != nil && ctx.Err() == nil {
			slog.Error("redis subscriber exited", "error", err)
		}
	}()

	m := models.New(Db)
	s := services.New(m, ws, redisClient)
	b := auth.NewTokenBucket(envInt("RATE_LIMIT_CAPACITY", 5), envInt("RATE_LIMIT_RATE", 2))
	c := controllers.New(s, b)

	r := routes.InitializeRoutes(c)

	r.Get("/health", Health)
	r.Get("/ready", Ready)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				slog.Info("metrics snapshot", "snapshot", metrics.Snapshot())
			case <-ctx.Done():
				return
			}
		}
	}()

	server := &http.Server{
		Addr:    ":" + envString("PORT", "8000"),
		Handler: r,
	}

	go func() {
		slog.Info("server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped unexpectedly", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hub.Cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		slog.Error("database pool close failed", "error", err)
		return
	}

	slog.Info("server exited cleanly")

}

func envInt(name string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil || value < 0 {
		return fallback
	}
	return value
}

func envString(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
