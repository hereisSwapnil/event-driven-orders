package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type createOrderRequest struct {
	ID            string `json:"id"`
	CustomerName  string `json:"customer_name"`
	TotalPrice    int64  `json:"total_price"`
	ScheduledTime string `json:"scheduled_time,omitempty"`
}

type orderResponse struct {
	ID            string `json:"id"`
	CustomerName  string `json:"customer_name"`
	TotalPrice    int64  `json:"total_price"`
	Status        string `json:"status"`
	ScheduledTime string `json:"scheduled_time,omitempty"`
	CreatedAt     string `json:"created_at"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var scheduledAt *time.Time
	if req.ScheduledTime != "" {
		t, err := time.Parse(time.RFC3339, req.ScheduledTime)
		if err != nil {
			http.Error(w, "invalid scheduled_time format", http.StatusBadRequest)
			return
		}
		scheduledAt = &t
	}

	order, err := h.service.CreateOrder(
		r.Context(),
		req.CustomerName,
		int(req.TotalPrice),
		scheduledAt,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := orderResponse{
		ID:           order.ID,
		CustomerName: order.CustomerName,
		TotalPrice:   int64(order.Price),
		Status:       string(order.Status),
		CreatedAt:    order.CreatedAt.Format(time.RFC3339),
	}

	if order.ScheduledAt != nil {
		resp.ScheduledTime = order.ScheduledAt.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	resp := orderResponse{
		ID:           order.ID,
		CustomerName: order.CustomerName,
		TotalPrice:   int64(order.Price),
		Status:       string(order.Status),
		CreatedAt:    order.CreatedAt.Format(time.RFC3339),
	}

	if order.ScheduledAt != nil {
		resp.ScheduledTime = order.ScheduledAt.Format(time.RFC3339)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
