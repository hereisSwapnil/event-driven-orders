package events

import (
	"context"

	"github.com/hereisSwapnil/event-driven-orders/internal/platform/kafka"
)


type OrderEventProducer struct {
	producer *kafka.Producer
}

func NewOrderEventProducer(producer *kafka.Producer) *OrderEventProducer {
	return &OrderEventProducer{
		producer: producer,
	}
}

func (p *OrderEventProducer) PublishOrderCreated(
	ctx context.Context,
	event OrderCreatedEvent,
) error {
	return p.producer.Publish(ctx, event.OrderID, event)
}