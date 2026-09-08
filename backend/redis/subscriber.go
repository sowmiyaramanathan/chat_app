package redis

import (
	"backend/entities/packet"
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

func (redis *redispubsub) SubscribeMessages(ctx context.Context, handler func(packet.MessageEvent)) error {
	for {
		err := redis.subscribeOnce(ctx, handler)
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			slog.Error("redis subscriber stopped; retrying", "error", err)
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
		}
	}
}

func (redis *redispubsub) subscribeOnce(ctx context.Context, handler func(packet.MessageEvent)) error {
	sub := redis.rdClient.Subscribe(ctx, messageChannel)
	defer sub.Close()

	if _, err := sub.Receive(ctx); err != nil {
		return err
	}
	slog.Info("redis subscriber connected", "channel", messageChannel)

	for {
		message, err := sub.ReceiveMessage(ctx)
		if err != nil {
			return err
		}

		var event packet.MessageEvent
		if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
			slog.Warn("invalid redis message event", "error", err)
			continue
		}
		if event.EventID == "" || event.RecipientID == "" || event.SenderID == "" {
			slog.Warn("redis message event missing routing fields")
			continue
		}

		handler(event)
	}
}
