package domain

import "time"

type OrderStatus string

const (
	OrderStatusCreated OrderStatus = "CREATED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusFailed OrderStatus = "FAILED"
)

type Order struct {
	ID string
	CustomerName string
	Price int
	Status OrderStatus
	ScheduledAt *time.Time
	CreatedAt time.Time
}