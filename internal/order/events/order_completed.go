package events

import "time"

type OrderCompletedEvent struct {
	OrderID     string    `json:"order_id"`
	CompletedAt time.Time `json:"completed_at"`
}
