package restaurant

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

var allowedStatusTransitions = map[model.OrderStatus]model.OrderStatus{
	model.OrderStatusAccepted:   model.OrderStatusPreparing,
	model.OrderStatusPreparing:  model.OrderStatusReady,
	model.OrderStatusReady:      model.OrderStatusDelivering,
	model.OrderStatusDelivering: model.OrderStatusDelivered,
}

func (s *Service) UpdateOrderStatus(
	ctx context.Context,
	input model.UpdateOrderStatus,
) (model.OrderStatusChange, error) {
	if input.RestaurantID == uuid.Nil || input.OrderID == uuid.Nil {
		return model.OrderStatusChange{}, errs.InvalidArgument("restaurant id and order id are required")
	}
	if !validOrderStatus(input.ExpectedStatus) || !validOrderStatus(input.NewStatus) {
		return model.OrderStatusChange{}, errs.InvalidArgument("order status is invalid")
	}

	var change model.OrderStatusChange
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		current, err := s.repository.GetOrderForUpdate(txCtx, input.OrderID)
		if err != nil {
			return err
		}

		change = model.OrderStatusChange{
			OrderID:   current.ID,
			Status:    current.Status,
			UpdatedAt: current.UpdatedAt,
		}
		if current.Status == input.NewStatus {
			return nil
		}
		if current.Status != input.ExpectedStatus {
			return errs.FailedPrecondition("order status differs from expected status")
		}
		if allowedStatusTransitions[current.Status] != input.NewStatus {
			return errs.FailedPrecondition("order status transition is not allowed")
		}

		updatedAt, err := s.repository.UpdateOrderStatus(txCtx, input.OrderID, input.NewStatus)
		if err != nil {
			return err
		}
		fromStatus := current.Status
		if err := s.repository.CreateOrderStatusHistory(txCtx, model.OrderStatusHistory{
			OrderID:    input.OrderID,
			FromStatus: &fromStatus,
			ToStatus:   input.NewStatus,
		}); err != nil {
			return err
		}

		change.Status = input.NewStatus
		change.UpdatedAt = updatedAt
		return nil
	})
	if err != nil {
		return model.OrderStatusChange{}, err
	}

	if err := s.kitchen.UpdateOrderStatus(ctx, input); err != nil {
		return model.OrderStatusChange{}, err
	}
	return change, nil
}

func validOrderStatus(status model.OrderStatus) bool {
	switch status {
	case model.OrderStatusAccepted,
		model.OrderStatusPreparing,
		model.OrderStatusReady,
		model.OrderStatusDelivering,
		model.OrderStatusDelivered:
		return true
	default:
		return false
	}
}
