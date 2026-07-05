-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
    order_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid UUID NOT NULL,
    part_uuids TEXT[] NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING_PAYMENT',
    transaction_uuid UUID,
    payment_method VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Индекс для быстрого поиска по user_uuid
CREATE INDEX idx_orders_user_uuid ON orders(user_uuid);

-- Индекс для быстрого поиска по статусу
CREATE INDEX idx_orders_status ON orders(status);

-- +goose Down
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_uuid;
DROP TABLE IF EXISTS orders;