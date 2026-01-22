package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func Ping(ctx context.Context, redisClient *redis.Client) error {
	return redisClient.Ping(ctx).Err()
}