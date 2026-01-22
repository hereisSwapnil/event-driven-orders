package events

import "time"

type OrderCreatedEvent struct {
	OrderID     string     `json:"order_id"`
	CreatedAt   time.Time  `json:"created_at"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
}
