-- +goose Up
UPDATE restaurants
SET accepting_orders = FALSE
WHERE is_active = FALSE AND accepting_orders = TRUE;

ALTER TABLE restaurants
    ADD CONSTRAINT restaurants_active_accepting_consistent
        CHECK (is_active OR NOT accepting_orders);

ALTER TABLE order_items
    DROP CONSTRAINT order_items_pkey,
    DROP CONSTRAINT order_items_order_menu_item_unique,
    DROP COLUMN id;

ALTER TABLE order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (order_id, menu_item_id);

ALTER TABLE order_status_history
    DROP COLUMN from_status;

ALTER TABLE order_status_history
    RENAME COLUMN to_status TO status;

-- +goose Down
ALTER TABLE order_status_history
    RENAME COLUMN status TO to_status;

ALTER TABLE order_status_history
    ADD COLUMN from_status order_status;

ALTER TABLE order_items
    DROP CONSTRAINT order_items_pkey,
    ADD COLUMN id UUID NOT NULL DEFAULT gen_random_uuid();

ALTER TABLE order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id),
    ADD CONSTRAINT order_items_order_menu_item_unique UNIQUE (order_id, menu_item_id);

ALTER TABLE restaurants
    DROP CONSTRAINT restaurants_active_accepting_consistent;
