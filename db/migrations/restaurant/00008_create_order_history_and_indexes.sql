-- +goose Up
CREATE TABLE order_status_history (
    id BIGSERIAL PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    from_status restaurant_order_status,
    to_status restaurant_order_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX orders_restaurant_status_created_idx
    ON orders (restaurant_id, status, created_at, id);

CREATE INDEX order_status_history_order_created_idx
    ON order_status_history (order_id, created_at, id);

-- +goose Down
DROP INDEX order_status_history_order_created_idx;
DROP INDEX orders_restaurant_status_created_idx;
DROP TABLE order_status_history;
