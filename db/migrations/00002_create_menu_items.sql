-- +goose Up
CREATE TABLE menu_items (
    id UUID PRIMARY KEY,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE RESTRICT,
    external_id VARCHAR(200) NOT NULL,
    name VARCHAR(200) NOT NULL,
    description VARCHAR(2000) NOT NULL DEFAULT '',
    price_minor BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT menu_items_restaurant_external_id_unique UNIQUE (restaurant_id, external_id),
    CONSTRAINT menu_items_external_id_not_blank CHECK (BTRIM(external_id) <> ''),
    CONSTRAINT menu_items_name_not_blank CHECK (BTRIM(name) <> ''),
    CONSTRAINT menu_items_price_non_negative CHECK (price_minor >= 0),
    CONSTRAINT menu_items_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT menu_items_version_positive CHECK (version > 0)
);

-- +goose Down
DROP TABLE menu_items;
