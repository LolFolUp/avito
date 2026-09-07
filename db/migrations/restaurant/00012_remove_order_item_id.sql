-- +goose Up
ALTER TABLE order_items
    DROP CONSTRAINT order_items_pkey,
    DROP CONSTRAINT order_items_order_menu_item_unique,
    DROP COLUMN id,
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (order_id, menu_item_id);

-- +goose Down
ALTER TABLE order_items
    DROP CONSTRAINT order_items_pkey,
    ADD COLUMN id UUID NOT NULL DEFAULT gen_random_uuid(),
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id),
    ADD CONSTRAINT order_items_order_menu_item_unique UNIQUE (order_id, menu_item_id);

ALTER TABLE order_items ALTER COLUMN id DROP DEFAULT;
