-- +goose Up
DROP INDEX menu_items_restaurant_available_idx;
DROP INDEX orders_restaurant_status_created_idx;

ALTER TABLE menu_items
    DROP CONSTRAINT menu_items_restaurant_id_fkey,
    DROP CONSTRAINT menu_items_restaurant_id_id_unique,
    DROP COLUMN restaurant_id;

ALTER TABLE orders
    DROP CONSTRAINT orders_restaurant_id_fkey,
    DROP COLUMN restaurant_id;

DROP TABLE restaurants;

CREATE INDEX menu_items_available_idx
    ON menu_items (id)
    WHERE is_active = TRUE AND is_available = TRUE;

CREATE INDEX orders_status_created_idx
    ON orders (status, created_at, id);

-- +goose Down
DROP INDEX orders_status_created_idx;
DROP INDEX menu_items_available_idx;

CREATE TABLE restaurants (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT restaurants_name_not_blank CHECK (BTRIM(name) <> ''),
    CONSTRAINT restaurants_currency_format CHECK (currency ~ '^[A-Z]{3}$')
);

INSERT INTO restaurants (id, name, currency)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Demo Kitchen',
    'RUB'
);

ALTER TABLE menu_items
    ADD COLUMN restaurant_id UUID NOT NULL
        DEFAULT '00000000-0000-0000-0000-000000000001'
        REFERENCES restaurants(id) ON DELETE RESTRICT,
    ADD CONSTRAINT menu_items_restaurant_id_id_unique UNIQUE (restaurant_id, id);

ALTER TABLE orders
    ADD COLUMN restaurant_id UUID NOT NULL
        DEFAULT '00000000-0000-0000-0000-000000000001'
        REFERENCES restaurants(id) ON DELETE RESTRICT;

CREATE INDEX menu_items_restaurant_available_idx
    ON menu_items (restaurant_id, id)
    WHERE is_active = TRUE AND is_available = TRUE;

CREATE INDEX orders_restaurant_status_created_idx
    ON orders (restaurant_id, status, created_at, id);
