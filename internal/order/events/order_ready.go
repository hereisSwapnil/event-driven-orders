package events

import "time"

type OrderReadyEvent struct {
	OrderID string `json:"order_id"`
	ReadyAt time.Time `json:"ready_at"`	
}
