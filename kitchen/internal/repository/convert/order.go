package convert

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/record"
)

func OrderToModel(source record.Order, items []record.OrderItem) model.Order {
	return model.Order{
		ID:              source.ID,
		RestaurantID:    source.RestaurantID,
		Status:          model.OrderStatus(source.Status),
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		Items:           OrderItemsToModel(items),
		TotalMinor:      source.TotalMinor,
		Currency:        source.Currency,
		RejectionReason: source.RejectionReason.String,
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
}

func OrderItemToModel(source record.OrderItem) model.OrderItem {
	return model.OrderItem{
		OrderID:        source.OrderID,
		MenuItemID:     source.MenuItemID,
		Name:           source.NameSnapshot,
		UnitPriceMinor: source.UnitPriceMinor,
		Quantity:       source.Quantity,
		LineTotalMinor: source.LineTotalMinor,
	}
}

func OrderItemsToModel(source []record.OrderItem) []model.OrderItem {
	result := make([]model.OrderItem, 0, len(source))
	for _, item := range source {
		result = append(result, OrderItemToModel(item))
	}

	return result
}
