package kitchen

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	convertrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/record"
)

const menuItemColumns = `
	id,
	restaurant_id,
	external_id,
	name,
	description,
	price_minor,
	currency,
	is_active,
	is_available,
	version,
	created_at,
	updated_at`

func (r *Repository) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error) {
	query := `SELECT ` + menuItemColumns + `
		FROM menu_items
		WHERE restaurant_id = $1 AND is_active = TRUE
		ORDER BY name, id`

	return r.queryMenuItems(ctx, query, restaurantID)
}

func (r *Repository) GetMenuForUpdate(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error) {
	query := `SELECT ` + menuItemColumns + `
		FROM menu_items
		WHERE restaurant_id = $1
		ORDER BY id
		FOR UPDATE`

	return r.queryMenuItems(ctx, query, restaurantID)
}

func (r *Repository) GetMenuItemsForUpdate(
	ctx context.Context,
	restaurantID uuid.UUID,
	itemIDs []uuid.UUID,
) ([]model.MenuItem, error) {
	query := `SELECT ` + menuItemColumns + `
		FROM menu_items
		WHERE restaurant_id = $1 AND id = ANY($2::uuid[])
		ORDER BY id
		FOR UPDATE`

	return r.queryMenuItems(ctx, query, restaurantID, itemIDs)
}

func (r *Repository) UpsertMenuItem(ctx context.Context, item model.MenuItem) error {
	query := `INSERT INTO menu_items (
			id,
			restaurant_id,
			external_id,
			name,
			description,
			price_minor,
			currency,
			is_active,
			is_available,
			version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (restaurant_id, external_id) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			price_minor = EXCLUDED.price_minor,
			currency = EXCLUDED.currency,
			is_active = EXCLUDED.is_active,
			is_available = EXCLUDED.is_available,
			version = EXCLUDED.version,
			updated_at = NOW()`

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx,
		query,
		item.ID,
		item.RestaurantID,
		item.ExternalID,
		item.Name,
		item.Description,
		item.PriceMinor,
		item.Currency,
		item.IsActive,
		item.IsAvailable,
		item.Version,
	)
	if err != nil {
		return repositoryError("upsert menu item", "menu item", err)
	}

	return nil
}

func (r *Repository) DeactivateMenuItemsExcept(
	ctx context.Context,
	restaurantID uuid.UUID,
	externalIDs []string,
) (int64, error) {
	query := `UPDATE menu_items
		SET is_active = FALSE,
			is_available = FALSE,
			version = version + 1,
			updated_at = NOW()
		WHERE restaurant_id = $1
			AND is_active = TRUE
			AND NOT (external_id = ANY($2::text[]))`

	result, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, restaurantID, externalIDs)
	if err != nil {
		return 0, repositoryError("deactivate missing menu items", "menu item", err)
	}

	return result.RowsAffected(), nil
}

func (r *Repository) queryMenuItems(ctx context.Context, query string, args ...any) ([]model.MenuItem, error) {
	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, args...)
	if err != nil {
		return nil, repositoryError("query menu items", "menu items", err)
	}
	defer rows.Close()

	items := make([]record.MenuItem, 0)
	for rows.Next() {
		var item record.MenuItem
		if err := rows.Scan(
			&item.ID,
			&item.RestaurantID,
			&item.ExternalID,
			&item.Name,
			&item.Description,
			&item.PriceMinor,
			&item.Currency,
			&item.IsActive,
			&item.IsAvailable,
			&item.Version,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, repositoryError("scan menu item", "menu item", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, repositoryError("iterate menu items", "menu items", err)
	}

	return convertrepo.MenuItemsToModel(items), nil
}
