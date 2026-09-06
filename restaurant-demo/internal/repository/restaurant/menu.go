package restaurant

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"
)

func (r *Repository) ListMenu(ctx context.Context) ([]model.MenuItem, error) {
	const query = `
		SELECT id, name, description, price_minor, currency,
		       is_active, is_available, available_quantity, version
		FROM menu_items
		WHERE is_active = TRUE
		ORDER BY name, id`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query)
	if err != nil {
		return nil, repositoryError("list menu", "menu", err)
	}
	defer rows.Close()

	items := make([]record.MenuItem, 0)

	for rows.Next() {
		var item record.MenuItem
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.PriceMinor,
			&item.Currency,
			&item.IsActive,
			&item.IsAvailable,
			&item.AvailableQuantity,
			&item.Version,
		)
		if err != nil {
			return nil, repositoryError("scan menu item", "menu item", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, repositoryError("iterate menu", "menu", err)
	}
	return convert.MenuItemsToModel(items), nil
}

func (r *Repository) GetMenuItemsForUpdate(
	ctx context.Context,
	itemIDs []uuid.UUID,
) ([]model.MenuItem, error) {
	const query = `
		SELECT id, name, description, price_minor, currency,
		       is_active, is_available, available_quantity, version
		FROM menu_items
		WHERE id = ANY($1)
		ORDER BY id
		FOR UPDATE`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, itemIDs)
	if err != nil {
		return nil, repositoryError("lock menu items", "menu items", err)
	}
	defer rows.Close()

	items := make([]record.MenuItem, 0, len(itemIDs))
	for rows.Next() {
		var item record.MenuItem
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.PriceMinor,
			&item.Currency,
			&item.IsActive,
			&item.IsAvailable,
			&item.AvailableQuantity,
			&item.Version,
		); err != nil {
			return nil, repositoryError("scan menu item", "menu item", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, repositoryError("iterate menu items", "menu items", err)
	}
	return convert.MenuItemsToModel(items), nil
}

func (r *Repository) ReserveMenuItem(
	ctx context.Context,
	itemID uuid.UUID,
	quantity int,
) error {
	const query = `
		UPDATE menu_items
		SET available_quantity = available_quantity - $2,
		    is_available = CASE WHEN available_quantity - $2 = 0 THEN FALSE ELSE is_available END,
		    version = CASE WHEN available_quantity - $2 = 0 THEN version + 1 ELSE version END,
		    updated_at = NOW()
		WHERE id = $1
		  AND is_active = TRUE
		  AND is_available = TRUE
		  AND available_quantity >= $2`

	result, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, itemID, quantity)
	if err != nil {
		return repositoryError("reserve menu item", "menu item", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("menu item %s was not reserved", itemID)
	}
	return nil
}
