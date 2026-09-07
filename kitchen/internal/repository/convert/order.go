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
		Status:          model.OrderStatus(source.Status),
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		Items:           OrderItemsToModel(items),
		TotalPrice:      source.TotalPrice,
		Currency:        source.Currency,
		RejectionReason: source.RejectionReason.String,
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
}
func OrderToRecord(model model.Order) record.Order {
	return record.Order{
		ID:              model.ID,
		RestaurantID:    model.RestaurantID,
		Status:          modelStatusToRecord(model.Status),
		CustomerName:    model.CustomerName,
		CustomerPhone:   model.CustomerPhone,
		DeliveryAddress: model.DeliveryAddress,
		TotalPrice:      model.TotalPrice,
		Currency:        model.Currency,
		RejectionReason: RejectionReasonToRecord(model.RejectionReason),
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func OrderItemToRecord(item []model.OrderItem) []record.OrderItem {
	items := make([]record.OrderItem, 0, len(item))
	for _, item := range item {
		items = append(items, orderItemToRecord(item))
	}
	return items
}

func orderItemToRecord(model model.OrderItem) record.OrderItem {
	return record.OrderItem{
		OrderID:        model.OrderID,
		MenuItemID:     model.MenuItemID,
		Name:           model.Name,
		UnitPriceMinor: model.UnitPriceMinor,
		Quantity:       model.Quantity,
		LineTotalPrice: model.LineTotalPrice,
	}
}

func RejectionReasonToRecord(reason string) pgtype.Text {
	return pgtype.Text{
		String: reason,
		Valid:  reason != "",
	}
}

func modelStatusToRecord(status model.OrderStatus) string {
	switch status {
	case model.OrderStatusPendingConfirmation:
		return "PendingConfirmation"
	case model.OrderStatusAccepted:
		return "Accepted"
	case model.OrderStatusPreparing:
		return "Preparing"
	case model.OrderStatusDelivering:
		return "Delivering"
	case model.OrderStatusRejected:
		return "Rejected"
	case model.OrderStatusDelivered:
		return "Delivered"
	default:
		return ""
	}
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
