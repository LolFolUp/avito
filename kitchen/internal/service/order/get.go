package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func (s *Service) Get(ctx context.Context, orderID uuid.UUID) (model.Order, error) {
	if orderID == uuid.Nil {
		return model.Order{}, errs.InvalidArgument("order id is required")
	}

	return s.orders.GetOrder(ctx, orderID)
}
