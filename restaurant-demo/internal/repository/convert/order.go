package convert

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/record"
)

func ConvertModelOrderToRecord(model model.Order) record.Order {
	return record.Order{
		ID:              model.ID,
		Status:          convertStatusToRecord(model.Status),
		CustomerName:    model.CustomerName,
		CustomerPhone:   model.CustomerPhone,
		DeliveryAddress: model.DeliveryAddress,
		TotalPrice:      model.TotalPrice,
		Currency:        model.Currency,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func ConvertModelOrderItemsToRecord(model []model.OrderItem) []record.OrderItem {
	records := make([]record.OrderItem, 0, len(model))
	for _, item := range model {
		rec := record.OrderItem{
			MenuItemID:     item.MenuItemID,
			Name:           item.Name,
			Price:          item.Price,
			Quantity:       item.Quantity,
			LineTotalPrice: item.LineTotalPrice,
		}
		records = append(records, rec)
	}
	return records
}
func convertStatusToRecord(status model.OrderStatus) string {
	switch status {
	case model.OrderStatusAccepted:
		return "ACCEPTED"
	case model.OrderStatusPreparing:
		return "PREPARING"
	case model.OrderStatusReady:
		return "READY"
	case model.OrderStatusDelivered:
		return "DELIVERED"
	case model.OrderStatusDelivering:
		return "DELIVERING"
	default:
		return ""
	}
}
