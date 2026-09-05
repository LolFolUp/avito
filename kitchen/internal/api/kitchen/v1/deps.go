package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

type CatalogService interface {
	ListRestaurants(ctx context.Context) ([]model.Restaurant, error)
	GetMenu(ctx context.Context, restaurantID uuid.UUID) (model.RestaurantMenu, error)
}

type OrderService interface {
	Create(ctx context.Context, input model.CreateOrder) (model.Order, error)
	Get(ctx context.Context, orderID uuid.UUID) (model.Order, error)
}

type ReadinessChecker interface {
	Ping(ctx context.Context) error
}
