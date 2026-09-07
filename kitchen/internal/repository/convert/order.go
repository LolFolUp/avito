package convert

import (
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/record"
)

func OrderToModel(source record.Order, items []record.OrderItem) model.Order {
	return model.Order{
		ID:              source.ID,
		RestaurantID:    source.RestaurantID,
		Status:          OrderStatusToModel(source.Status),
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		Items:           OrderItemsToModel(items),
		TotalPrice:      source.TotalPrice,
		Currency:        source.Currency,
		RejectionReason: RejectionReasonToModel(source.RejectionReason),
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
}
func OrderToRecord(source model.Order) record.Order {
	return record.Order{
		ID:              source.ID,
		RestaurantID:    source.RestaurantID,
		Status:          OrderStatusToRecord(source.Status),
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		TotalPrice:      source.TotalPrice,
		Currency:        source.Currency,
		RejectionReason: RejectionReasonToRecord(source.RejectionReason),
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
}

func OrderItemsToRecords(source []model.OrderItem) []record.OrderItem {
	items := make([]record.OrderItem, 0, len(source))
	for _, item := range source {
		items = append(items, orderItemToRecord(item))
	}
	return items
}

func orderItemToRecord(source model.OrderItem) record.OrderItem {
	return record.OrderItem{
		OrderID:        source.OrderID,
		MenuItemID:     source.MenuItemID,
		Name:           source.Name,
		UnitPriceMinor: source.UnitPriceMinor,
		Quantity:       source.Quantity,
		LineTotalPrice: source.LineTotalPrice,
	}
}

func RejectionReasonToRecord(reason string) pgtype.Text {
	return pgtype.Text{
		String: reason,
		Valid:  reason != "",
	}
}

func RejectionReasonToModel(reason pgtype.Text) string {
	if !reason.Valid {
		return ""
	}
	return reason.String
}

func OrderStatusToRecord(status model.OrderStatus) string {
	return string(status)
}

func OrderStatusToModel(status string) model.OrderStatus {
	return model.OrderStatus(status)
}

func OrderItemToModel(source record.OrderItem) model.OrderItem {
	return model.OrderItem{
		OrderID:        source.OrderID,
		MenuItemID:     source.MenuItemID,
		Name:           source.Name,
		UnitPriceMinor: source.UnitPriceMinor,
		Quantity:       source.Quantity,
		LineTotalPrice: source.LineTotalPrice,
	}
}

func OrderItemsToModel(source []record.OrderItem) []model.OrderItem {
	result := make([]model.OrderItem, 0, len(source))
	for _, item := range source {
		result = append(result, OrderItemToModel(item))
	}

	return result
}
