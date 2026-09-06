package convert

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"
)

func convertMenuItemToModel(record record.MenuItem) model.MenuItem {
	return model.MenuItem{
		ID:                record.ID,
		Name:              record.Name,
		Description:       record.Description,
		PriceMinor:        record.PriceMinor,
		Currency:          record.Currency,
		IsActive:          record.IsActive,
		IsAvailable:       record.IsAvailable,
		AvailableQuantity: record.AvailableQuantity,
		Version:           record.Version,
	}
}

func MenuItemsToModel(records []record.MenuItem) []model.MenuItem {
	items := make([]model.MenuItem, len(records))
	for i, recordItem := range records {
		items[i] = convertMenuItemToModel(recordItem)
	}
	return items
}
