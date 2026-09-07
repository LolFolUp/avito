package restaurant

import (
	"context"
	"github.com/google/uuid"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

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
		&rec.ID,
		&rec.Status,
		&rec.CustomerName,
		&rec.CustomerPhone,
		&rec.DeliveryAddress,
		&rec.TotalPrice,
		&rec.Currency,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		return repositoryError("create order", "order", err)
	}
	return r.createOrderItems(ctx, order.ID, recItems)
}

func (r *Repository) createOrderItems(ctx context.Context, orderID uuid.UUID, items []record.OrderItem) error {
	// TODO проверь код еще раз
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
	// TODO проверь код еще раз
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
