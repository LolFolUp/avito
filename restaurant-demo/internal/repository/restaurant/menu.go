package restaurant

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"
)

func (r *Repository) GetRestaurant(ctx context.Context, id uuid.UUID) (model.Restaurant, error) {
	const query = `
		SELECT id, name, currency
		FROM restaurants
		WHERE id = $1`

	var restaurant record.Restaurant
	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Currency,
	)
	if err != nil {
		return model.Restaurant{}, repositoryError("get restaurant", "restaurant", err)
	}
	return convert.RestaurantToModel(restaurant), nil
}

func (r *Repository) ListMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error) {
	const query = `
		SELECT id, restaurant_id, name, description, price_minor, currency,
		       is_active, is_available, available_quantity, version
		FROM menu_items
		WHERE restaurant_id = $1 AND is_active = TRUE
		ORDER BY name, id`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, restaurantID)
	if err != nil {
		return nil, repositoryError("list menu", "menu", err)
	}
	defer rows.Close()

	items := make([]record.MenuItem, 0)
	for rows.Next() {
		item, err := scanMenuItem(rows)
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
	restaurantID uuid.UUID,
	itemIDs []uuid.UUID,
) ([]model.MenuItem, error) {
	const query = `
		SELECT id, restaurant_id, name, description, price_minor, currency,
		       is_active, is_available, available_quantity, version
		FROM menu_items
		WHERE restaurant_id = $1 AND id = ANY($2)
		ORDER BY id
		FOR UPDATE`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, restaurantID, itemIDs)
	if err != nil {
		return nil, repositoryError("lock menu items", "menu items", err)
	}
	defer rows.Close()

	items := make([]record.MenuItem, 0, len(itemIDs))
	for rows.Next() {
		item, err := scanMenuItem(rows)
		if err != nil {
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
	restaurantID uuid.UUID,
	itemID uuid.UUID,
	quantity int,
) error {
	const query = `
		UPDATE menu_items
		SET available_quantity = available_quantity - $3,
		    is_available = CASE WHEN available_quantity - $3 = 0 THEN FALSE ELSE is_available END,
		    version = CASE WHEN available_quantity - $3 = 0 THEN version + 1 ELSE version END,
		    updated_at = NOW()
		WHERE restaurant_id = $1
		  AND id = $2
		  AND is_active = TRUE
		  AND is_available = TRUE
		  AND available_quantity >= $3`

	result, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, restaurantID, itemID, quantity)
	if err != nil {
		return repositoryError("reserve menu item", "menu item", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("menu item %s was not reserved", itemID)
	}
	return nil
}

type menuItemScanner interface {
	Scan(dest ...any) error
}

func scanMenuItem(scanner menuItemScanner) (record.MenuItem, error) {
	var item record.MenuItem
	err := scanner.Scan(
		&item.ID,
		&item.RestaurantID,
		&item.Name,
		&item.Description,
		&item.PriceMinor,
		&item.Currency,
		&item.IsActive,
		&item.IsAvailable,
		&item.AvailableQuantity,
		&item.Version,
	)
	return item, err
}
