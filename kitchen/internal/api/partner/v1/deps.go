package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

type OrderService interface {
	UpdateStatus(
		ctx context.Context,
		restaurantID uuid.UUID,
		orderID uuid.UUID,
		expectedStatus model.OrderStatus,
		newStatus model.OrderStatus,
		reason string,
	) (model.Order, error)
}
