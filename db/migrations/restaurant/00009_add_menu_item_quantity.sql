-- +goose Up
ALTER TABLE menu_items
    ADD COLUMN available_quantity INTEGER NOT NULL DEFAULT 100,
    ADD CONSTRAINT menu_items_available_quantity_non_negative CHECK (available_quantity >= 0);

-- +goose Down
ALTER TABLE menu_items
    DROP CONSTRAINT menu_items_available_quantity_non_negative,
    DROP COLUMN available_quantity;
