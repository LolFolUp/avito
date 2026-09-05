package convert

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/record"
)

func MenuItemToModel(source record.MenuItem) model.MenuItem {
	return model.MenuItem{
		ID:           source.ID,
		RestaurantID: source.RestaurantID,
		ExternalID:   source.ExternalID,
		Name:         source.Name,
		Description:  source.Description,
		PriceMinor:   source.PriceMinor,
		Currency:     source.Currency,
		IsActive:     source.IsActive,
		IsAvailable:  source.IsAvailable,
		Version:      source.Version,
		CreatedAt:    source.CreatedAt,
		UpdatedAt:    source.UpdatedAt,
	}
}

func MenuItemsToModel(source []record.MenuItem) []model.MenuItem {
	result := make([]model.MenuItem, 0, len(source))
	for _, item := range source {
		result = append(result, MenuItemToModel(item))
	}

	return result
}
