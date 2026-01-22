package repository

import (
	"context"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/domain"
)

type OrderCache interface {
	Get(ctx context.Context, id string) (*domain.Order, error)
	Set(ctx context.Context, order *domain.Order) error
}
