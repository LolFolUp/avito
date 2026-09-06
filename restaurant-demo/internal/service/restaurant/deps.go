package restaurant

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type RestaurantRepository interface {
	GetRestaurant(ctx context.Context, id uuid.UUID) (model.Restaurant, error)
	ListMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error)
	GetMenuItemsForUpdate(ctx context.Context, restaurantID uuid.UUID, itemIDs []uuid.UUID) ([]model.MenuItem, error)
	ReserveMenuItem(ctx context.Context, restaurantID uuid.UUID, itemID uuid.UUID, quantity int) error
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
	CreateOrderItems(ctx context.Context, orderID uuid.UUID, items []model.OrderItem) error
	CreateOrderStatusHistory(ctx context.Context, history model.OrderStatusHistory) error
	GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error)
}
