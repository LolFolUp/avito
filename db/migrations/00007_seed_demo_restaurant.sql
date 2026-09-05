-- +goose Up
INSERT INTO restaurants (
    id,
    name,
    is_active,
    accepting_orders,
    currency
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Demo Kitchen',
    TRUE,
    TRUE,
    'RUB'
);

-- +goose Down
DELETE FROM restaurants
WHERE id = '00000000-0000-0000-0000-000000000001';
