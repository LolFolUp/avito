package catalog

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

type RestaurantRepository interface {
	ListRestaurants(ctx context.Context) ([]model.Restaurant, error)
	GetRestaurant(ctx context.Context, id uuid.UUID) (model.Restaurant, error)
	GetRestaurantForUpdate(ctx context.Context, id uuid.UUID) (model.Restaurant, error)
}

type MenuRepository interface {
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error)
	GetMenuForUpdate(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error)
	UpsertMenuItem(ctx context.Context, item model.MenuItem) error
	DeactivateMenuItemsExcept(ctx context.Context, restaurantID uuid.UUID, externalIDs []string) (int64, error)
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
