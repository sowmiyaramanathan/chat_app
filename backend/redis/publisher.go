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

	receivers, err := redis.rdClient.Publish(ctx, messageChannel, data).Result()
	if err != nil {
		return err
	}
	slog.Info("redis message published", "channel", messageChannel, "event_id", event.EventID, "receivers", receivers)
	return nil
}
