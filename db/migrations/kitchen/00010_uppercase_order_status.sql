-- +goose Up
ALTER TYPE order_status RENAME VALUE 'pending_confirmation' TO 'PENDING_CONFIRMATION';
ALTER TYPE order_status RENAME VALUE 'accepted' TO 'ACCEPTED';
ALTER TYPE order_status RENAME VALUE 'preparing' TO 'PREPARING';
ALTER TYPE order_status RENAME VALUE 'ready' TO 'READY';
ALTER TYPE order_status RENAME VALUE 'delivering' TO 'DELIVERING';
ALTER TYPE order_status RENAME VALUE 'delivered' TO 'DELIVERED';
ALTER TYPE order_status RENAME VALUE 'rejected' TO 'REJECTED';

-- +goose Down
ALTER TYPE order_status RENAME VALUE 'PENDING_CONFIRMATION' TO 'pending_confirmation';
ALTER TYPE order_status RENAME VALUE 'ACCEPTED' TO 'accepted';
ALTER TYPE order_status RENAME VALUE 'PREPARING' TO 'preparing';
ALTER TYPE order_status RENAME VALUE 'READY' TO 'ready';
ALTER TYPE order_status RENAME VALUE 'DELIVERING' TO 'delivering';
ALTER TYPE order_status RENAME VALUE 'DELIVERED' TO 'delivered';
ALTER TYPE order_status RENAME VALUE 'REJECTED' TO 'rejected';
