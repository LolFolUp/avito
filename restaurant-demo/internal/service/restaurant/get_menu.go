package restaurant

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (s *Service) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, string, error) {
	if restaurantID == uuid.Nil {
		return nil, "", errs.InvalidArgument("restaurant id is required")
	}
	if restaurantID != s.restaurantID {
		return nil, "", errs.NotFound("restaurant not found")
	}

	items, err := s.repository.ListMenu(ctx)
	if err != nil {
		return nil, "", err
	}
	return items, s.currency, nil
}
