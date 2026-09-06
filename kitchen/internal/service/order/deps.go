package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

type RestaurantRepository interface {
	GetRestaurant(ctx context.Context, id uuid.UUID) (model.Restaurant, error)
}

type RestaurantClient interface {
	CreateOrder(ctx context.Context, orderID uuid.UUID, input model.CreateOrder) (model.Order, error)
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	CreateOrderItems(ctx context.Context, items []model.OrderItem) error
	CreateOrderStatusHistory(ctx context.Context, history model.OrderStatusHistory) error
	GetOrder(ctx context.Context, id uuid.UUID) (model.Order, error)
	GetOrderForUpdate(ctx context.Context, id uuid.UUID) (model.Order, error)
	ListOrders(ctx context.Context, filter model.OrderFilter) ([]model.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus, reason string) error
}

type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
