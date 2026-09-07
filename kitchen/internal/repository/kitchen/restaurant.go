package kitchen

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	convertrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/convert"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/record"
)

const restaurantColumns = `
	id,
	name,
	is_active,
	accepting_orders,
	currency,
	created_at,
	updated_at`

func (r *Repository) ListRestaurants(ctx context.Context) ([]model.Restaurant, error) {
	query := `SELECT ` + restaurantColumns + `
		FROM restaurants
		WHERE is_active = TRUE
		ORDER BY name, id`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query)
	if err != nil {
		return nil, repositoryError("list restaurants", "restaurants", err)
	}
	defer rows.Close()

	restaurants := make([]record.Restaurant, 0)
	for rows.Next() {
		var restaurant record.Restaurant
		if err := rows.Scan(
			&restaurant.ID,
			&restaurant.Name,
			&restaurant.IsActive,
			&restaurant.AcceptingOrders,
			&restaurant.Currency,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		); err != nil {
			return nil, repositoryError("scan restaurant", "restaurant", err)
		}
		restaurants = append(restaurants, restaurant)
	}
	if err := rows.Err(); err != nil {
		return nil, repositoryError("iterate restaurants", "restaurants", err)
	}

	return convertrepo.RestaurantsToModel(restaurants), nil
}

func (r *Repository) GetRestaurant(ctx context.Context, id uuid.UUID) (model.Restaurant, error) {
	query := `SELECT ` + restaurantColumns + ` FROM restaurants WHERE id = $1`

	var restaurant record.Restaurant
	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.IsActive,
		&restaurant.AcceptingOrders,
		&restaurant.Currency,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)
	if err != nil {
		return model.Restaurant{}, repositoryError("get restaurant", "restaurant", err)
	}

	return convertrepo.RestaurantToModel(restaurant), nil
}
