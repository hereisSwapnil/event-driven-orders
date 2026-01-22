package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/events"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/handler"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/repository/postgres"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/repository/redis"
	"github.com/hereisSwapnil/event-driven-orders/internal/order/service"
	"github.com/hereisSwapnil/event-driven-orders/internal/platform/kafka"
	redisClient "github.com/hereisSwapnil/event-driven-orders/internal/platform/redis"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/order_service?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redisClient.NewRedisClient("localhost:6379")
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	kafkaProducer := kafka.NewProducer("localhost:9092", "order.created")
	orderEventProducer := events.NewOrderEventProducer(kafkaProducer)


	repo := postgres.NewOrderRepository(db)
	cache := redis.NewOrderCache(redisClient, 5*time.Minute)

	orderService := service.NewOrderService(repo, cache, orderEventProducer)
	orderHandler := handler.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /orders", orderHandler.GetOrder)

	log.Println("order service running on :8080")
	http.ListenAndServe(":8080", mux)
}
