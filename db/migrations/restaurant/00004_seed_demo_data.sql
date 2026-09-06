-- +goose Up
INSERT INTO restaurants (
    id,
    name,
    currency
) VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Demo Kitchen',
    'RUB'
);

INSERT INTO menu_items (
    restaurant_id,
    external_id,
    name,
    description,
    price_minor,
    currency,
    is_active,
    is_available
) VALUES
    (
        '00000000-0000-0000-0000-000000000001',
        'pizza-margherita',
        'Margherita',
        'Tomato sauce and mozzarella',
        59000,
        'RUB',
        TRUE,
        TRUE
    ),
    (
        '00000000-0000-0000-0000-000000000001',
        'pasta-carbonara',
        'Carbonara',
        'Pasta, bacon and cream sauce',
        65000,
        'RUB',
        TRUE,
        TRUE
    ),
    (
        '00000000-0000-0000-0000-000000000001',
        'tiramisu',
        'Tiramisu',
        'Classic Italian dessert',
        35000,
        'RUB',
        TRUE,
        FALSE
    );

-- +goose Down
DELETE FROM menu_items
WHERE restaurant_id = '00000000-0000-0000-0000-000000000001';

DELETE FROM restaurants
WHERE id = '00000000-0000-0000-0000-000000000001';
