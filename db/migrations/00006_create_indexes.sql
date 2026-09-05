-- +goose Up
CREATE INDEX restaurants_active_idx
    ON restaurants (id)
    WHERE is_active = TRUE;

CREATE INDEX menu_items_restaurant_active_idx
    ON menu_items (restaurant_id, is_active, is_available);

CREATE INDEX orders_restaurant_status_created_idx
    ON orders (restaurant_id, status, created_at, id);

CREATE INDEX order_status_history_order_created_idx
    ON order_status_history (order_id, created_at, id);

-- +goose Down
DROP INDEX order_status_history_order_created_idx;
DROP INDEX orders_restaurant_status_created_idx;
DROP INDEX menu_items_restaurant_active_idx;
DROP INDEX restaurants_active_idx;
