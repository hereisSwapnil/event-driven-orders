CREATE TABLE orders (
    id TEXT PRIMARY KEY,
    customer_name TEXT NOT NULL,
    price BIGINT NOT NULL,
    status TEXT NOT NULL,
    scheduled_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_orders_status ON orders(status);