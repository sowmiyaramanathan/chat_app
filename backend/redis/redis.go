package redis

import (
	"backend/entities/packet"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type redispubsub struct {
	rdClient *redis.Client
}

type RedisPubSub interface {
	ConnectionRegistry
	PublishMessage(ctx context.Context, event packet.MessageEvent) (err error)
	SubscribeMessages(ctx context.Context, handler func(packet.MessageEvent)) error
	Close() error
}

func New(ctx context.Context) (RedisPubSub, error) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_HOST")
	}
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := rdClient.Ping(pingCtx).Err(); err != nil {
		_ = rdClient.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return &redispubsub{rdClient: rdClient}, nil
}

func (redis *redispubsub) Close() error {
	return redis.rdClient.Close()
}
