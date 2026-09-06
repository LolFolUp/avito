package restaurant

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (s *Service) GetMenu(ctx context.Context, restaurantID uuid.UUID) (model.Restaurant, []model.MenuItem, error) {
	if restaurantID == uuid.Nil {
		return model.Restaurant{}, nil, errs.InvalidArgument("restaurant id is required")
	}

	restaurant, err := s.repository.GetRestaurant(ctx, restaurantID)
	if err != nil {
		return model.Restaurant{}, nil, err
	}
	items, err := s.repository.ListMenu(ctx, restaurantID)
	if err != nil {
		return model.Restaurant{}, nil, err
	}
	return restaurant, items, nil
}
