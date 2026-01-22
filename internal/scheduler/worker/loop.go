package worker

import (
	"context"
	"log"
	"time"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/events"
	"github.com/hereisSwapnil/event-driven-orders/internal/platform/kafka"
)

type Loop struct {
	scheduler *RedisScheduler
	producer  *kafka.Producer
}

func NewLoop(
	scheduler *RedisScheduler,
	producer *kafka.Producer,
) *Loop {
	return &Loop{
		scheduler: scheduler,
		producer:  producer,
	}
}

func (l *Loop) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now().UTC()

			orderIDs, err := l.scheduler.PopDue(ctx, now)
			if err != nil {
				log.Println("scheduler error:", err)
				continue
			}

			for _, orderID := range orderIDs {
				event := events.OrderReadyEvent{
					OrderID: orderID,
					ReadyAt: now,
				}

				err := l.producer.Publish(ctx, orderID, event)
				if err != nil {
					log.Println("publish error:", err)
					continue
				}

				_ = l.scheduler.Remove(ctx, orderID)
				log.Println("order ready:", orderID)
			}
		case <-ctx.Done():
			return
		}
	}
}
