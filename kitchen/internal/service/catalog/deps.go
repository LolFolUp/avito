package catalog

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

type RestaurantRepository interface {
	ListRestaurants(ctx context.Context) ([]model.Restaurant, error)
	GetRestaurant(ctx context.Context, id uuid.UUID) (model.Restaurant, error)
}

type RestaurantClient interface {
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error)
}
