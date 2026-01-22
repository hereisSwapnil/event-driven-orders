package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/events"
	"github.com/segmentio/kafka-go"
)

type OrderCreatedConsumer struct {
	reader    *kafka.Reader
	scheduler *RedisScheduler
}

func NewOrderCreatedConsumer(
	brokers []string,
	topic string,
	groupID string,
	scheduler *RedisScheduler,
) *OrderCreatedConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &OrderCreatedConsumer{
		reader:    reader,
		scheduler: scheduler,
	}
}

func (c *OrderCreatedConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		var event events.OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid event:", err)
			continue
		}

		log.Printf("order created: %s scheduled_at: %v", event.OrderID, event.ScheduledAt)

		// Calculate execution time. If no scheduled time, execute immediately.
		executeAt := time.Now().UTC()
		if event.ScheduledAt != nil {
			executeAt = *event.ScheduledAt
		}

		if err := c.scheduler.Add(ctx, event.OrderID, executeAt); err != nil {
			log.Println("failed to schedule order:", err)
			continue
		}

		log.Printf("scheduled order %s at %v", event.OrderID, executeAt)
	}
}
