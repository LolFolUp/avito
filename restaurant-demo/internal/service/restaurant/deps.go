package restaurant

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

type PartnerClient interface {
	SyncMenu(ctx context.Context, restaurantID uuid.UUID, menu []model.MenuItem) (model.MenuSyncResult, error)
	ListOrders(
		ctx context.Context,
		restaurantID uuid.UUID,
		status model.OrderStatus,
		pageSize int,
		pageToken string,
	) (model.OrderPage, error)
	UpdateOrderStatus(
		ctx context.Context,
		restaurantID uuid.UUID,
		orderID uuid.UUID,
		expectedStatus model.OrderStatus,
		newStatus model.OrderStatus,
		reason string,
	) (model.Order, error)
}

type MenuRepository interface {
	ListMenu(ctx context.Context) ([]model.MenuItem, error)
	ItemsAvailable(ctx context.Context, externalIDs []string) (bool, error)
}
