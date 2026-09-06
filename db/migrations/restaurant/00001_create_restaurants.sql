-- +goose Up
CREATE TABLE restaurants (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT restaurants_name_not_blank CHECK (BTRIM(name) <> ''),
    CONSTRAINT restaurants_currency_format CHECK (currency ~ '^[A-Z]{3}$')
);

-- +goose Down
DROP TABLE restaurants;
