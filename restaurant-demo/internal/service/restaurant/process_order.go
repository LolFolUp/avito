package restaurant

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (s *Service) processOrder(ctx context.Context, order model.Order) error {
	externalIDs := make([]string, 0, len(order.Items))
	for _, item := range order.Items {
		externalIDs = append(externalIDs, item.ExternalItemID)
	}

	available, err := s.menu.ItemsAvailable(ctx, externalIDs)
	if err != nil {
		return err
	}

	status := model.OrderStatusAccepted
	reason := ""
	if !available {
		status = model.OrderStatusRejected
		reason = "one or more menu items are unavailable"
	}

	updated, err := s.partner.UpdateOrderStatus(
		ctx,
		s.restaurantID,
		order.ID,
		model.OrderStatusPendingConfirmation,
		status,
		reason,
	)
	if err != nil {
		return err
	}

	s.logger.Info("order processed", "order_id", order.ID, "status", updated.Status)
	return nil
}
