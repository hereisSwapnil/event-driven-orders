package worker

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisScheduler struct {
	client *redis.Client
	key    string
}

func NewRedisScheduler(client *redis.Client) *RedisScheduler {
	return &RedisScheduler{
		client: client,
		key:    "scheduled_orders",
	}
}

func (r *RedisScheduler) Add(ctx context.Context, orderID string, executeAt time.Time) error {
	score := float64(executeAt.Unix())
	return r.client.ZAdd(ctx, r.key, redis.Z{
		Score:  score,
		Member: orderID,
	}).Err()
}

func (r *RedisScheduler) PopDue(ctx context.Context, now time.Time) ([]string, error) {
	max := strconv.FormatInt(now.Unix(), 10)

	return r.client.ZRangeByScore(ctx, r.key, &redis.ZRangeBy{
		Min: "0",
		Max: max,
	}).Result()
}

func (r *RedisScheduler) Remove(ctx context.Context, orderID string) error {
	return r.client.ZRem(ctx, r.key, orderID).Err()
}
