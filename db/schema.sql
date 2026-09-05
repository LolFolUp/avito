CREATE TABLE restaurants (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    accepting_orders BOOLEAN NOT NULL DEFAULT TRUE,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT restaurants_name_not_blank CHECK (BTRIM(name) <> ''),
    CONSTRAINT restaurants_currency_format CHECK (currency ~ '^[A-Z]{3}$')
);

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

CREATE TYPE order_status AS ENUM (
    'pending_confirmation',
    'accepted',
    'preparing',
    'ready',
    'delivering',
    'delivered',
    'rejected'
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE RESTRICT,
    status order_status NOT NULL,
    customer_name VARCHAR(200) NOT NULL,
    customer_phone VARCHAR(32) NOT NULL,
    delivery_address VARCHAR(1000) NOT NULL,
    total_minor BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    rejection_reason VARCHAR(1000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT orders_customer_name_not_blank CHECK (BTRIM(customer_name) <> ''),
    CONSTRAINT orders_customer_phone_not_blank CHECK (BTRIM(customer_phone) <> ''),
    CONSTRAINT orders_delivery_address_not_blank CHECK (BTRIM(delivery_address) <> ''),
    CONSTRAINT orders_total_non_negative CHECK (total_minor >= 0),
    CONSTRAINT orders_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT orders_rejection_reason_consistent CHECK (
        (status = 'rejected' AND NULLIF(BTRIM(rejection_reason), '') IS NOT NULL)
        OR (status <> 'rejected' AND rejection_reason IS NULL)
    )
);

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

CREATE TABLE order_status_history (
    id BIGSERIAL PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    from_status order_status,
    to_status order_status NOT NULL,
    reason VARCHAR(1000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX restaurants_active_idx
    ON restaurants (id)
    WHERE is_active = TRUE;

CREATE INDEX menu_items_restaurant_active_idx
    ON menu_items (restaurant_id, is_active, is_available);

CREATE INDEX orders_restaurant_status_created_idx
    ON orders (restaurant_id, status, created_at, id);

CREATE INDEX order_status_history_order_created_idx
    ON order_status_history (order_id, created_at, id);
