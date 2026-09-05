package restaurant

import (
	"context"
	"time"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

const pageSize = 50

func (s *Service) Run(ctx context.Context, interval time.Duration) error {
	if err := s.SyncMenu(ctx); err != nil {
		return err
	}
	if err := s.pollOnce(ctx); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := s.pollOnce(ctx); err != nil {
				s.logger.Error("poll orders", "error", err)
			}
		}
	}
}

func (s *Service) pollOnce(ctx context.Context) error {
	pageToken := ""
	for {
		page, err := s.partner.ListOrders(
			ctx,
			s.restaurantID,
			model.OrderStatusPendingConfirmation,
			pageSize,
			pageToken,
		)
		if err != nil {
			return err
		}

		for _, order := range page.Orders {
			if err := s.processOrder(ctx, order); err != nil {
				s.logger.Error("process order", "order_id", order.ID, "error", err)
			}
		}

		if page.NextPageToken == "" {
			return nil
		}
		pageToken = page.NextPageToken
	}
}
