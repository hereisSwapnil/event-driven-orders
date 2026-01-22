package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/events"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/service"
)

type OrderReadyConsumer struct {
	reader       *kafka.Reader
	orderService *service.OrderService
}

func NewOrderReadyConsumer(
	brokers []string,
	topic string,
	groupID string,
	orderService *service.OrderService,
) *OrderReadyConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	return &OrderReadyConsumer{
		reader:       reader,
		orderService: orderService,
	}
}

func (c *OrderReadyConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		var event events.OrderReadyEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("invalid event:", err)
			continue
		}

		log.Println("processing order:", event.OrderID)

		err = c.orderService.MarkOrderAsProcessing(ctx, event.OrderID)
		if err != nil {
			log.Println("failed to mark processing:", err)
			continue
		}

		time.Sleep(2 * time.Second)

		err = c.orderService.CompleteOrder(ctx, event.OrderID)

		if err != nil {
			log.Println("failed to mark completed:", err)
			continue
		}

		log.Println("order completed:", event.OrderID)
	}
}
