package restaurant

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

const orderColumns = `
	id,
	status,
	customer_name,
	customer_phone,
	delivery_address,
	total_minor,
	currency,
	created_at,
	updated_at`

func (r *Repository) CreateOrder(ctx context.Context, order model.Order) error {
	const query = `
		INSERT INTO orders (
			id, status, customer_name, customer_phone,
			delivery_address, total_minor, currency, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	rec := convert.ConvertModelOrderToRecord(order)
	recItems := convert.ConvertModelOrderItemsToRecord(order.Items)

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx,
		query,
		rec.ID,
		rec.Status,
		rec.CustomerName,
		rec.CustomerPhone,
		rec.DeliveryAddress,
		rec.TotalPrice,
		rec.Currency,
		rec.CreatedAt,
		rec.UpdatedAt,
	)
	if err != nil {
		return repositoryError("create order", "order", err)
	}
	return r.createOrderItems(ctx, order.ID, recItems)
}

func (r *Repository) createOrderItems(ctx context.Context, orderID uuid.UUID, items []record.OrderItem) error {
	const query = `
		INSERT INTO order_items (
			order_id, menu_item_id, name_snapshot, unit_price_minor, quantity
		) VALUES ($1, $2, $3, $4, $5)`

	executor := r.getter.DefaultTrOrDB(ctx, r.pool)

	for _, item := range items {
		_, err := executor.Exec(
			ctx,
			query,
			orderID,
			item.MenuItemID,
			item.Name,
			item.Price,
			item.Quantity,
		)
		if err != nil {
			return repositoryError("create order item", "order item", err)
		}
	}
	return nil
}

func (r *Repository) CreateOrderStatusHistory(ctx context.Context, history model.OrderStatusHistory) error {
	const query = `
		INSERT INTO order_status_history (order_id, from_status, to_status)
		VALUES ($1, $2, $3)`

	var fromStatus any
	if history.FromStatus != nil {
		fromStatus = *history.FromStatus
	}
	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx,
		query,
		history.OrderID,
		fromStatus,
		history.ToStatus,
	)
	if err != nil {
		return repositoryError("create order status history", "order status history", err)
	}
	return nil
}

func (r *Repository) GetOrderForUpdate(ctx context.Context, id uuid.UUID) (model.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders WHERE id = $1 FOR UPDATE`

	var order record.Order
	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.Status,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.DeliveryAddress,
		&order.TotalPrice,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, repositoryError("get order for update", "order", err)
	}
	return convert.OrderToModel(order), nil
}

func (r *Repository) UpdateOrderStatus(
	ctx context.Context,
	id uuid.UUID,
	status model.OrderStatus,
) (time.Time, error) {
	const query = `
		UPDATE orders
		SET status = $2,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	var updatedAt time.Time
	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(
		ctx,
		query,
		id,
		convert.OrderStatusToRecord(status),
	).Scan(&updatedAt)
	if err != nil {
		return time.Time{}, repositoryError("update order status", "order", err)
	}
	return updatedAt, nil
}
