package catalog

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func (s *Service) ListRestaurants(ctx context.Context) ([]model.Restaurant, error) {
	return s.restaurants.ListRestaurants(ctx)
}

func (s *Service) GetMenu(ctx context.Context, restaurantID uuid.UUID) (model.RestaurantMenu, error) {
	if restaurantID == uuid.Nil {
		return model.RestaurantMenu{}, errs.InvalidArgument("restaurant id is required")
	}

	restaurant, err := s.restaurants.GetRestaurant(ctx, restaurantID)
	if err != nil {
		return model.RestaurantMenu{}, err
	}
	if !restaurant.IsActive {
		return model.RestaurantMenu{}, errs.NotFound("restaurant")
	}

	items, err := s.partner.GetMenu(ctx, restaurantID)
	if err != nil {
		return model.RestaurantMenu{}, err
	}

	return model.RestaurantMenu{Restaurant: restaurant, Items: items}, nil
}
