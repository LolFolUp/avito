package v1

import (
	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func menuToProto(source []model.MenuItem) []*partnerpb.MenuItem {
	result := make([]*partnerpb.MenuItem, 0, len(source))
	for _, item := range source {
		result = append(result, &partnerpb.MenuItem{
			ExternalId:  item.ExternalID,
			Name:        item.Name,
			Description: item.Description,
			PriceMinor:  item.PriceMinor,
			Currency:    item.Currency,
			IsActive:    item.IsActive,
			IsAvailable: item.IsAvailable,
		})
	}
	return result
}

func orderFromProto(source *partnerpb.Order) (model.Order, error) {
	orderID, err := uuid.Parse(source.GetId())
	if err != nil {
		return model.Order{}, err
	}
	restaurantID, err := uuid.Parse(source.GetRestaurantId())
	if err != nil {
		return model.Order{}, err
	}

	items := make([]model.OrderItem, 0, len(source.GetItems()))
	for _, item := range source.GetItems() {
		menuItemID, err := uuid.Parse(item.GetMenuItemId())
		if err != nil {
			return model.Order{}, err
		}
		items = append(items, model.OrderItem{
			MenuItemID:     menuItemID,
			ExternalItemID: item.GetExternalItemId(),
			Name:           item.GetName(),
			UnitPriceMinor: item.GetUnitPriceMinor(),
			Quantity:       int(item.GetQuantity()),
			LineTotalMinor: item.GetLineTotalMinor(),
		})
	}

	result := model.Order{
		ID:              orderID,
		RestaurantID:    restaurantID,
		Status:          statusFromProto(source.GetStatus()),
		CustomerName:    source.GetCustomerName(),
		CustomerPhone:   source.GetCustomerPhone(),
		DeliveryAddress: source.GetDeliveryAddress(),
		Items:           items,
		TotalMinor:      source.GetTotalMinor(),
		Currency:        source.GetCurrency(),
		RejectionReason: source.GetRejectionReason(),
	}
	if source.GetCreatedAt() != nil {
		result.CreatedAt = source.GetCreatedAt().AsTime()
	}
	if source.GetUpdatedAt() != nil {
		result.UpdatedAt = source.GetUpdatedAt().AsTime()
	}
	return result, nil
}

func statusFromProto(source partnerpb.OrderStatus) model.OrderStatus {
	switch source {
	case partnerpb.OrderStatus_ORDER_STATUS_PENDING_CONFIRMATION:
		return model.OrderStatusPendingConfirmation
	case partnerpb.OrderStatus_ORDER_STATUS_ACCEPTED:
		return model.OrderStatusAccepted
	case partnerpb.OrderStatus_ORDER_STATUS_PREPARING:
		return model.OrderStatusPreparing
	case partnerpb.OrderStatus_ORDER_STATUS_READY:
		return model.OrderStatusReady
	case partnerpb.OrderStatus_ORDER_STATUS_DELIVERING:
		return model.OrderStatusDelivering
	case partnerpb.OrderStatus_ORDER_STATUS_DELIVERED:
		return model.OrderStatusDelivered
	case partnerpb.OrderStatus_ORDER_STATUS_REJECTED:
		return model.OrderStatusRejected
	default:
		return ""
	}
}

func statusToProto(source model.OrderStatus) partnerpb.OrderStatus {
	switch source {
	case model.OrderStatusPendingConfirmation:
		return partnerpb.OrderStatus_ORDER_STATUS_PENDING_CONFIRMATION
	case model.OrderStatusAccepted:
		return partnerpb.OrderStatus_ORDER_STATUS_ACCEPTED
	case model.OrderStatusPreparing:
		return partnerpb.OrderStatus_ORDER_STATUS_PREPARING
	case model.OrderStatusReady:
		return partnerpb.OrderStatus_ORDER_STATUS_READY
	case model.OrderStatusDelivering:
		return partnerpb.OrderStatus_ORDER_STATUS_DELIVERING
	case model.OrderStatusDelivered:
		return partnerpb.OrderStatus_ORDER_STATUS_DELIVERED
	case model.OrderStatusRejected:
		return partnerpb.OrderStatus_ORDER_STATUS_REJECTED
	default:
		return partnerpb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
