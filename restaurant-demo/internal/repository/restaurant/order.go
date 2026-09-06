package restaurant

import (
	"context"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (r *Repository) CreateOrder(ctx context.Context, order model.Order) error {
	const query = `
		INSERT INTO orders (
			id, status, customer_name, customer_phone,
			delivery_address, total_minor, currency
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	var rec record.Order
	res, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx,
		query,
		rec.ID,
		rec.Status,
		rec.CustomerName,
		rec.CustomerPhone,
		rec.DeliveryAddress,
		rec.TotalMinor,
		rec.Currency,
	)
	if err != nil {
		return repositoryError("create order", "order", err)
	}
	if res.RowsAffected() == 0 {
		return repositoryError("create order", "order", err)
	}
	return nil
}

func (r *Repository) CreateOrderItems(ctx context.Context, orderID uuid.UUID, items []model.OrderItem) error {
	const query = `
		INSERT INTO order_items (
			id, order_id, menu_item_id, name_snapshot, unit_price_minor, quantity
		) VALUES ($1, $2, $3, $4, $5, $6)`

	executor := r.getter.DefaultTrOrDB(ctx, r.pool)
	for _, item := range items {
		_, err := executor.Exec(
			ctx,
			query,
			uuid.New(),
			orderID,
			item.MenuItemID,
			item.Name,
			item.UnitPriceMinor,
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
