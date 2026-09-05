-- +goose Up
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

-- +goose Down
DROP TABLE orders;
DROP TYPE order_status;
