package kitchen

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	convertrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/record"
)

const orderColumns = `
	id,
	restaurant_id,
	status,
	customer_name,
	customer_phone,
	delivery_address,
	total_minor,
	currency,
	rejection_reason,
	created_at,
	updated_at`

const orderItemColumns = `
	order_id,
	menu_item_id,
	name_snapshot,
	unit_price_minor,
	quantity,
	line_total_minor`

func (r *Repository) CreateOrder(ctx context.Context, order model.Order) error {
	query := `INSERT INTO orders (
			id,
			restaurant_id,
			status,
			customer_name,
			customer_phone,
			delivery_address,
			total_minor,
			currency,
			rejection_reason
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULL)`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
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
	)
	if err != nil {
		return repositoryError("create order", "order", err)
	}

	return nil
}

func (r *Repository) CreateOrderItems(ctx context.Context, items []model.OrderItem) error {
	query := `INSERT INTO order_items (
			order_id,
			menu_item_id,
			name_snapshot,
			unit_price_minor,
			quantity
		) VALUES ($1, $2, $3, $4, $5)`

	executor := r.getter.DefaultTrOrDB(ctx, r.pool)
	for _, item := range items {
		_, err := executor.Exec(
			ctx,
			query,
			item.OrderID,
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
	query := `INSERT INTO order_status_history (
			order_id,
			status,
			reason
		) VALUES ($1, $2, $3)`

	var reason any
	if history.Reason != "" {
		reason = history.Reason
	}

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx,
		query,
		history.OrderID,
		history.Status,
		reason,
	)
	if err != nil {
		return repositoryError("create order status history", "order status history", err)
	}

	return nil
}

func (r *Repository) GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, err := r.getOrderRecord(ctx, id, false)
	if err != nil {
		return model.Order{}, err
	}

	items, err := r.getOrderItemRecords(ctx, id)
	if err != nil {
		return model.Order{}, err
	}

	return convertrepo.OrderToModel(order, items), nil
}

func (r *Repository) GetOrderForUpdate(ctx context.Context, id uuid.UUID) (model.Order, error) {
	order, err := r.getOrderRecord(ctx, id, true)
	if err != nil {
		return model.Order{}, err
	}

	return convertrepo.OrderToModel(order, nil), nil
}

func (r *Repository) ListOrders(ctx context.Context, filter model.OrderFilter) ([]model.Order, error) {
	query := `SELECT ` + orderColumns + `
		FROM orders
		WHERE restaurant_id = $1
			AND status = $2
			AND (
				$3::timestamptz IS NULL
				OR (created_at, id) > ($3::timestamptz, $4::uuid)
			)
		ORDER BY created_at, id
		LIMIT $5`

	var afterCreatedAt any
	var afterID any
	if filter.AfterCreatedAt != nil && filter.AfterID != nil {
		afterCreatedAt = *filter.AfterCreatedAt
		afterID = *filter.AfterID
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(
		ctx,
		query,
		filter.RestaurantID,
		filter.Status,
		afterCreatedAt,
		afterID,
		filter.Limit,
	)
	if err != nil {
		return nil, repositoryError("list orders", "orders", err)
	}
	defer rows.Close()

	orderRecords := make([]record.Order, 0)
	for rows.Next() {
		order, err := scanOrder(rows.Scan)
		if err != nil {
			return nil, repositoryError("scan order", "order", err)
		}
		orderRecords = append(orderRecords, order)
	}
	if err := rows.Err(); err != nil {
		return nil, repositoryError("iterate orders", "orders", err)
	}

	orders := make([]model.Order, 0, len(orderRecords))
	for _, orderRecord := range orderRecords {
		items, err := r.getOrderItemRecords(ctx, orderRecord.ID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, convertrepo.OrderToModel(orderRecord, items))
	}

	return orders, nil
}

func (r *Repository) UpdateOrderStatus(
	ctx context.Context,
	id uuid.UUID,
	status model.OrderStatus,
	reason string,
) error {
	query := `UPDATE orders
		SET status = $2,
			rejection_reason = NULLIF($3, ''),
			updated_at = NOW()
		WHERE id = $1`

	result, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, id, status, reason)
	if err != nil {
		return repositoryError("update order status", "order", err)
	}
	if result.RowsAffected() == 0 {
		return errs.NotFound("order")
	}

	return nil
}

func (r *Repository) getOrderRecord(ctx context.Context, id uuid.UUID, forUpdate bool) (record.Order, error) {
	query := `SELECT ` + orderColumns + ` FROM orders WHERE id = $1`
	if forUpdate {
		query += ` FOR UPDATE`
	}

	order, err := scanOrder(func(dest ...any) error {
		return r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, id).Scan(dest...)
	})
	if err != nil {
		return record.Order{}, repositoryError("get order", "order", err)
	}

	return order, nil
}

func (r *Repository) getOrderItemRecords(ctx context.Context, orderID uuid.UUID) ([]record.OrderItem, error) {
	query := `SELECT ` + orderItemColumns + `
		FROM order_items
		WHERE order_id = $1
		ORDER BY menu_item_id`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, orderID)
	if err != nil {
		return nil, repositoryError("get order items", "order items", err)
	}
	defer rows.Close()

	items := make([]record.OrderItem, 0)
	for rows.Next() {
		var item record.OrderItem
		if err := rows.Scan(
			&item.OrderID,
			&item.MenuItemID,
			&item.NameSnapshot,
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

func scanOrder(scan func(dest ...any) error) (record.Order, error) {
	var order record.Order
	err := scan(
		&order.ID,
		&order.RestaurantID,
		&order.Status,
		&order.CustomerName,
		&order.CustomerPhone,
		&order.DeliveryAddress,
		&order.TotalMinor,
		&order.Currency,
		&order.RejectionReason,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	return order, err
}
