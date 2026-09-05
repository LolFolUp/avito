package convert

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/record"
)

func RestaurantToModel(source record.Restaurant) model.Restaurant {
	return model.Restaurant{
		ID:              source.ID,
		Name:            source.Name,
		IsActive:        source.IsActive,
		AcceptingOrders: source.AcceptingOrders,
		Currency:        source.Currency,
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
}

func RestaurantsToModel(source []record.Restaurant) []model.Restaurant {
	result := make([]model.Restaurant, 0, len(source))
	for _, restaurant := range source {
		result = append(result, RestaurantToModel(restaurant))
	}

	return result
}
