package order

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func (s *Service) UpdateStatus(
	ctx context.Context,
	restaurantID uuid.UUID,
	orderID uuid.UUID,
	expectedStatus model.OrderStatus,
	newStatus model.OrderStatus,
	reason string,
) (model.Order, error) {
	reason = strings.TrimSpace(reason)
	if restaurantID == uuid.Nil || orderID == uuid.Nil {
		return model.Order{}, errs.InvalidArgument("restaurant id and order id are required")
	}
	if !validStatus(expectedStatus) || !validStatus(newStatus) {
		return model.Order{}, errs.InvalidArgument("order status is invalid")
	}
	if newStatus == model.OrderStatusRejected && (reason == "" || len(reason) > 1000) {
		return model.Order{}, errs.InvalidArgument("rejection reason is required")
	}
	if newStatus != model.OrderStatusRejected && reason != "" {
		return model.Order{}, errs.InvalidArgument("reason is only allowed when rejecting an order")
	}

	var updated model.Order
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		current, err := s.orders.GetOrderForUpdate(txCtx, orderID)
		if err != nil {
			return err
		}
		if current.RestaurantID != restaurantID {
			return errs.NotFound("order")
		}
		if current.Status != expectedStatus {
			return errs.FailedPrecondition(
				errs.CodeStatusConflict,
				"order status differs from expected status",
				nil,
			)
		}
		if current.Status == newStatus {
			updated, err = s.orders.GetOrder(txCtx, orderID)
			return err
		}
		if !canTransition(current.Status, newStatus) {
			return errs.FailedPrecondition(
				errs.CodeInvalidStatusTransition,
				"order status transition is not allowed",
				nil,
			)
		}

		if err := s.orders.UpdateOrderStatus(txCtx, orderID, newStatus, reason); err != nil {
			return err
		}
		fromStatus := current.Status
		if err := s.orders.CreateOrderStatusHistory(txCtx, model.OrderStatusHistory{
			OrderID:    orderID,
			FromStatus: &fromStatus,
			ToStatus:   newStatus,
			Reason:     reason,
		}); err != nil {
			return err
		}

		updated, err = s.orders.GetOrder(txCtx, orderID)
		return err
	})
	if err != nil {
		return model.Order{}, err
	}

	return updated, nil
}
