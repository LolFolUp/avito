package restaurant

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type RestaurantRepository interface {
	ListMenu(ctx context.Context) ([]model.MenuItem, error)
	GetMenuItemsForUpdate(ctx context.Context, itemIDs []uuid.UUID) ([]model.MenuItem, error)
	ReserveMenuItem(ctx context.Context, itemID uuid.UUID, quantity int) error
	CreateOrder(ctx context.Context, order model.Order) error
	CreateOrderStatusHistory(ctx context.Context, history model.OrderStatusHistory) error
	GetOrderForUpdate(ctx context.Context, id uuid.UUID) (model.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) (time.Time, error)
}

type KitchenClient interface {
	UpdateOrderStatus(ctx context.Context, input model.UpdateOrderStatus) error
}
