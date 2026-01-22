package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/events"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/repository/postgres"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/service"
	"github.com/hereisSwapnil/event-driven-orders/internal/platform/kafka"
	"github.com/hereisSwapnil/event-driven-orders/internal/processing/consumer"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewOrderRepository(db)

	kafkaProducer := kafka.NewProducer(
		strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		"order.completed",
	)

	orderEventProducer := events.NewOrderEventProducer(kafkaProducer)

	orderService := service.NewOrderService(repo, nil, orderEventProducer)

	readyConsumer := consumer.NewOrderReadyConsumer(
		strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		"order.ready",
		"processing-group",
		orderService,
	)

	log.Println("processing service started")
	readyConsumer.Start(ctx)
}
