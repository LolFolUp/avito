-- +goose Up
CREATE TABLE order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
    external_item_id VARCHAR(200) NOT NULL,
    name_snapshot VARCHAR(200) NOT NULL,
    unit_price_minor BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    line_total_minor BIGINT GENERATED ALWAYS AS (unit_price_minor * quantity) STORED,
    CONSTRAINT order_items_order_menu_item_unique UNIQUE (order_id, menu_item_id),
    CONSTRAINT order_items_external_id_not_blank CHECK (BTRIM(external_item_id) <> ''),
    CONSTRAINT order_items_name_not_blank CHECK (BTRIM(name_snapshot) <> ''),
    CONSTRAINT order_items_unit_price_non_negative CHECK (unit_price_minor >= 0),
    CONSTRAINT order_items_quantity_range CHECK (quantity BETWEEN 1 AND 99)
);

-- +goose Down
DROP TABLE order_items;
