package convert

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"
)

func OrderToModel(source record.Order, items []model.OrderItem) model.Order {
	return model.Order{
		ID:              source.ID,
		RestaurantID:    source.RestaurantID,
		Status:          model.OrderStatus(source.Status),
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		Items:           items,
		TotalMinor:      source.TotalMinor,
		Currency:        source.Currency,
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
}
