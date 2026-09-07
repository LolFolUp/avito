package v1

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func menuItemsFromProto(restaurantID uuid.UUID, items []*partnerpb.MenuItem) ([]model.MenuItem, error) {
	result := make([]model.MenuItem, 0, len(items))
	for _, source := range items {
		id, err := uuid.Parse(source.GetId())
		if err != nil {
			return nil, fmt.Errorf("parse restaurant menu item id: %w", err)
		}
		result = append(result, model.MenuItem{
			ID:           id,
			RestaurantID: restaurantID,
			Name:         source.GetName(),
			Description:  source.GetDescription(),
			PriceMinor:   source.GetPriceMinor(),
			Currency:     source.GetCurrency(),
			IsAvailable:  source.GetIsAvailable(),
			Version:      source.GetVersion(),
		})
	}
	return result, nil
}

func createOrderRequestToProto(orderID uuid.UUID, input model.CreateOrder) *partnerpb.CreateRestaurantOrderRequest {
	items := make([]*partnerpb.CreateRestaurantOrderItem, 0, len(input.Items))
	for _, item := range input.Items {
		items = append(items, &partnerpb.CreateRestaurantOrderItem{
			MenuItemId:      item.MenuItemID.String(),
			Quantity:        int32(item.Quantity),
			ObservedVersion: item.ObservedVersion,
		})
	}
	return &partnerpb.CreateRestaurantOrderRequest{
		OrderId:         orderID.String(),
		RestaurantId:    input.RestaurantID.String(),
		CustomerName:    input.CustomerName,
		CustomerPhone:   input.CustomerPhone,
		DeliveryAddress: input.DeliveryAddress,
		Items:           items,
	}
}

func orderFromProto(source *partnerpb.RestaurantOrder) (model.Order, error) {
	if source == nil {
		return model.Order{}, fmt.Errorf("restaurant returned an empty order")
	}
	orderID, err := uuid.Parse(source.GetId())
	if err != nil {
		return model.Order{}, fmt.Errorf("parse restaurant order id: %w", err)
	}
	restaurantID, err := uuid.Parse(source.GetRestaurantId())
	if err != nil {
		return model.Order{}, fmt.Errorf("parse restaurant id: %w", err)
	}

	items := make([]model.OrderItem, 0, len(source.GetItems()))
	for _, sourceItem := range source.GetItems() {
		itemID, err := uuid.Parse(sourceItem.GetMenuItemId())
		if err != nil {
			return model.Order{}, fmt.Errorf("parse restaurant order item id: %w", err)
		}
		items = append(items, model.OrderItem{
			OrderID:        orderID,
			MenuItemID:     itemID,
			Name:           sourceItem.GetName(),
			UnitPriceMinor: sourceItem.GetUnitPriceMinor(),
			Quantity:       int(sourceItem.GetQuantity()),
			LineTotalPrice: sourceItem.GetLineTotalMinor(),
		})
	}

	return model.Order{
		ID:              orderID,
		RestaurantID:    restaurantID,
		Status:          orderStatusFromProto(source.GetStatus()),
		CustomerName:    source.GetCustomerName(),
		CustomerPhone:   source.GetCustomerPhone(),
		DeliveryAddress: source.GetDeliveryAddress(),
		Items:           items,
		TotalPrice:      source.GetTotalMinor(),
		Currency:        source.GetCurrency(),
		CreatedAt:       source.GetCreatedAt().AsTime(),
		UpdatedAt:       source.GetUpdatedAt().AsTime(),
	}, nil
}

func orderStatusFromProto(status partnerpb.RestaurantOrderStatus) model.OrderStatus {
	switch status {
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_ACCEPTED:
		return model.OrderStatusAccepted
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_PREPARING:
		return model.OrderStatusPreparing
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_READY:
		return model.OrderStatusReady
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_DELIVERING:
		return model.OrderStatusDelivering
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_DELIVERED:
		return model.OrderStatusDelivered
	default:
		return ""
	}
}
