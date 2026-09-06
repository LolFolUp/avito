-- +goose Up
CREATE INDEX menu_items_restaurant_available_idx
    ON menu_items (restaurant_id, external_id)
    WHERE is_active = TRUE AND is_available = TRUE;

-- +goose Down
DROP INDEX menu_items_restaurant_available_idx;
