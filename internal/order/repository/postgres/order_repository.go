package postgres

import (
	"context"
	"database/sql"

	"github.com/hereisSwapnil/event-driven-orders/internal/order/domain"
)

type OrderRepository struct {
	db *sql.DB	
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, customer_name, price, status, scheduled_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, order.ID, order.CustomerName, order.Price, order.Status, order.ScheduledAt, order.CreatedAt)

	return err
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `
		SELECT id, customer_name, price, status, scheduled_at, created_at
		FROM orders
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var order domain.Order
	err := row.Scan(&order.ID, &order.CustomerName, &order.Price, &order.Status, &order.ScheduledAt, &order.CreatedAt)

	return &order, err
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $2
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id, status)

	return err
}