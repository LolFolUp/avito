package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

type RestaurantService interface {
	GetMenu(ctx context.Context, restaurantID uuid.UUID) (model.Restaurant, []model.MenuItem, error)
	CreateOrder(ctx context.Context, input model.CreateOrder) (model.Order, error)
}
