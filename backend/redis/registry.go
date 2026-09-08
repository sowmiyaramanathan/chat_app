package redis

import (
	"context"
	"fmt"
	"time"
)

const connectionTTL = 30 * time.Second

type ConnectionRegistry interface {
	RegisterConnection(ctx context.Context, userID, connectionID, instanceID string) error
	RefreshConnection(ctx context.Context, userID, connectionID string) error
	UnregisterConnection(ctx context.Context, userID, connectionID, instanceID string) error
}

func (redis *redispubsub) RegisterConnection(ctx context.Context, userID, connectionID, instanceID string) error {
	key := connectionKey(userID, connectionID)
	pipe := redis.rdClient.TxPipeline()
	pipe.Set(ctx, key, instanceID, connectionTTL)
	pipe.SAdd(ctx, instanceConnectionsKey(instanceID), key)
	_, err := pipe.Exec(ctx)
	return err
}

func (redis *redispubsub) RefreshConnection(ctx context.Context, userID, connectionID string) error {
	return redis.rdClient.Expire(ctx, connectionKey(userID, connectionID), connectionTTL).Err()
}

func (redis *redispubsub) UnregisterConnection(ctx context.Context, userID, connectionID, instanceID string) error {
	pipe := redis.rdClient.TxPipeline()
	pipe.Del(ctx, connectionKey(userID, connectionID))
	pipe.SRem(ctx, instanceConnectionsKey(instanceID), connectionKey(userID, connectionID))
	_, err := pipe.Exec(ctx)
	return err
}

func connectionKey(userID, connectionID string) string {
	return fmt.Sprintf("ws:user:%s:connection:%s", userID, connectionID)
}

func instanceConnectionsKey(instanceID string) string {
	return fmt.Sprintf("ws:instance:%s:connections", instanceID)
}
