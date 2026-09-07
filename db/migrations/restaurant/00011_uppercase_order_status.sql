-- +goose Up
ALTER TYPE restaurant_order_status RENAME VALUE 'accepted' TO 'ACCEPTED';
ALTER TYPE restaurant_order_status RENAME VALUE 'preparing' TO 'PREPARING';
ALTER TYPE restaurant_order_status RENAME VALUE 'ready' TO 'READY';
ALTER TYPE restaurant_order_status RENAME VALUE 'delivering' TO 'DELIVERING';
ALTER TYPE restaurant_order_status RENAME VALUE 'delivered' TO 'DELIVERED';

-- +goose Down
ALTER TYPE restaurant_order_status RENAME VALUE 'ACCEPTED' TO 'accepted';
ALTER TYPE restaurant_order_status RENAME VALUE 'PREPARING' TO 'preparing';
ALTER TYPE restaurant_order_status RENAME VALUE 'READY' TO 'ready';
ALTER TYPE restaurant_order_status RENAME VALUE 'DELIVERING' TO 'delivering';
ALTER TYPE restaurant_order_status RENAME VALUE 'DELIVERED' TO 'delivered';
