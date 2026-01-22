package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
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
	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	// Run simple migration
	migration, err := os.ReadFile("migrations/postgres/001_create_orders.sql")
	if err != nil {
		log.Fatal("failed to read migration:", err)
	}

	for i := 0; i < 30; i++ {
		if _, err := db.Exec(string(migration)); err == nil {
			log.Println("migration ran successfully")
			break
		} else {
			log.Printf("failed to run migration: %v, retrying in 1s...", err)
			time.Sleep(1 * time.Second)
		}
		if i == 29 {
			log.Fatal("failed to run migration after retries")
		}
	}

	redisCli := redisClient.NewRedisClient(os.Getenv("REDIS_ADDR"))
	if err := redisCli.Ping(context.Background()).Err(); err != nil {
		log.Printf("failed to connect to redis: %v", err)
	}

	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	for i := 0; i < 30; i++ {
		if err := kafka.CreateTopic(kafkaBrokers, "order.created"); err == nil {
			log.Println("kafka topic created/verified")
			break
		} else {
			log.Printf("failed to create topic: %v, retrying in 1s...", err)
			time.Sleep(1 * time.Second)
		}
	}

	kafkaProducer := kafka.NewProducer(
		kafkaBrokers,
		"order.created",
	)

	orderEventProducer := events.NewOrderEventProducer(kafkaProducer)

	repo := postgres.NewOrderRepository(db)
	cache := redis.NewOrderCache(redisCli, 5*time.Minute)

	orderService := service.NewOrderService(repo, cache, orderEventProducer)
	orderHandler := handler.NewOrderHandler(orderService)

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodGet:
			orderHandler.GetOrder(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("order service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
