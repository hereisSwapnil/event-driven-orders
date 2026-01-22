package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/hereisSwapnil/event-driven-orders/internal/platform/kafka"
	"github.com/hereisSwapnil/event-driven-orders/internal/scheduler/worker"
)

func main() {
	ctx := context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	scheduler := worker.NewRedisScheduler(redisClient)

	producer := kafka.NewProducer(
		strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		"order.ready",
	)

	consumer := worker.NewOrderCreatedConsumer(
		strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		"order.created",
		"scheduler-group",
		scheduler,
	)

	loop := worker.NewLoop(scheduler, producer)

	log.Println("scheduler worker started")
	
	go consumer.Start(ctx)
	loop.Start(ctx)
}
