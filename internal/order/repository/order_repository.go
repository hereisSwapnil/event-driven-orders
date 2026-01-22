package repository

import (
	"context"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/domain"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) error
}