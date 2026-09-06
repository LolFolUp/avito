package restaurant

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"
)

func (r *Repository) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	const query = `
		INSERT INTO orders (
			id, restaurant_id, status, customer_name, customer_phone,
			delivery_address, total_minor, currency
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at`

	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(
		ctx,
		query,
		order.ID,
		order.RestaurantID,
		order.Status,
		order.CustomerName,
		order.CustomerPhone,
		order.DeliveryAddress,
		order.TotalMinor,
		order.Currency,
	).Scan(&order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return model.Order{}, repositoryError("create order", "order", err)
	}
	return order, nil
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

func (r *Repository) GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error) {
	const query = `
		SELECT id, restaurant_id, status, customer_name, customer_phone,
		       delivery_address, total_minor, currency, created_at, updated_at
		FROM orders
		WHERE id = $1`

	var order record.Order
	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.RestaurantID,
		&order.Status,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.DeliveryAddress,
		&order.TotalMinor,
		&order.Currency,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, repositoryError("get order", "order", err)
	}

	items, err := r.getOrderItems(ctx, id)
	if err != nil {
		return model.Order{}, err
	}
	return convert.OrderToModel(order, items), nil
}

func (r *Repository) getOrderItems(ctx context.Context, orderID uuid.UUID) ([]model.OrderItem, error) {
	const query = `
		SELECT menu_item_id, name_snapshot, unit_price_minor, quantity, line_total_minor
		FROM order_items
		WHERE order_id = $1
		ORDER BY id`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, orderID)
	if err != nil {
		return nil, repositoryError("get order items", "order items", err)
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0)
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.MenuItemID,
			&item.Name,
			&item.UnitPriceMinor,
			&item.Quantity,
			&item.LineTotalMinor,
		); err != nil {
			return nil, repositoryError("scan order item", "order item", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, repositoryError("iterate order items", "order items", err)
	}
	return items, nil
}
