package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/domain"
	"github.com/redis/go-redis/v9"
)

type OrderCache struct {
	client *redis.Client
	ttl time.Duration
}

func NewOrderCache(client *redis.Client, ttl time.Duration) *OrderCache {
	return &OrderCache{
		client: client,
		ttl: ttl,
	}
}

func (c *OrderCache) Get(ctx context.Context, id string) (*domain.Order, error) {
	key := "Order:" + id

	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var order domain.Order
	err = json.Unmarshal([]byte(data), &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (c *OrderCache) Set(ctx context.Context, order *domain.Order) error {
	key := "Order:" + order.ID

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, c.ttl).Err()
}
