package redis

import (
	"backend/entities/packet"
	"context"
	"encoding/json"
	"log/slog"
)

const messageChannel = "chat:messages"

func (redis *redispubsub) PublishMessage(ctx context.Context, event packet.MessageEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("error publishing message", "error", err)
		return err
	}

	return redis.rdClient.Publish(ctx, messageChannel, data).Err()
}
