package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/domain"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/events"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/repository"
)

type OrderService struct {
	repo      repository.OrderRepository
	cache     repository.OrderCache
	eventProd *events.OrderEventProducer
}


func NewOrderService(
	repo repository.OrderRepository,
	cache repository.OrderCache,
	eventProd *events.OrderEventProducer,
) *OrderService {
	return &OrderService{
		repo:      repo,
		cache:     cache,
		eventProd: eventProd,
	}
}


func (s *OrderService) CreateOrder(ctx context.Context, customerName string, totalPrice int, scheduledAt *time.Time) (*domain.Order, error) {
	if customerName == "" {
		return nil, errors.New("customer name is required")
	}

	if totalPrice <= 0 {
		return nil, errors.New("total price must be greater than 0")
	}

	order := &domain.Order{
		ID: uuid.New().String(),
		CustomerName: customerName,
		Price: totalPrice,
		Status: domain.OrderStatusCreated,
		ScheduledAt: scheduledAt,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, order); err != nil {
		return nil, err
	}

	if err := s.eventProd.PublishOrderCreated(ctx, events.OrderCreatedEvent{
		OrderID:     order.ID,
		CreatedAt:   order.CreatedAt,
		ScheduledAt: order.ScheduledAt,
	}); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	if order, err := s.cache.Get(ctx, id); err == nil {
		return order, nil
	}

	return s.repo.GetOrderByID(ctx, id)
}

func (s *OrderService) MarkOrderAsCompleted(ctx context.Context, id string) error {
	return s.repo.UpdateOrderStatus(ctx, id, domain.OrderStatusCompleted)
}

func (s *OrderService) MarkOrderAsFailed(ctx context.Context, id string) error {
	return s.repo.UpdateOrderStatus(ctx, id, domain.OrderStatusFailed)
}

func (s *OrderService) MarkOrderAsProcessing(ctx context.Context, id string) error {
	return s.repo.UpdateOrderStatus(ctx, id, domain.OrderStatusProcessing)
}

func (s *OrderService) CompleteOrder(
	ctx context.Context,
	orderID string,
) error {

	// 1. Update database status
	err := s.repo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusCompleted)
	if err != nil {
		return err
	}

	// 2. Publish completion event (best effort)
	if s.eventProd != nil {
		_ = s.eventProd.PublishOrderCompleted(
			ctx,
			events.OrderCompletedEvent{
				OrderID:     orderID,
				CompletedAt: time.Now().UTC(),
			},
		)
	}

	return nil
}
