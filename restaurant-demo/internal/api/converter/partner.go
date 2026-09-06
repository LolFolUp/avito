package converter

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func MenuItemsToProto(items []model.MenuItem) []*partnerpb.MenuItem {
	result := make([]*partnerpb.MenuItem, 0, len(items))
	for _, item := range items {
		result = append(result, &partnerpb.MenuItem{
			Id:          item.ID.String(),
			Name:        item.Name,
			Description: item.Description,
			PriceMinor:  item.PriceMinor,
			Currency:    item.Currency,
			IsAvailable: item.IsAvailable && item.AvailableQuantity > 0,
			Version:     item.Version,
		})
	}
	return result
}

func CreateOrderFromProto(request *partnerpb.CreateRestaurantOrderRequest) (model.CreateOrder, error) {
	orderID, err := uuid.Parse(request.GetOrderId())
	if err != nil {
		return model.CreateOrder{}, fmt.Errorf("invalid order id: %w", err)
	}
	restaurantID, err := uuid.Parse(request.GetRestaurantId())
	if err != nil {
		return model.CreateOrder{}, fmt.Errorf("invalid restaurant id: %w", err)
	}

	items := make([]model.CreateOrderItem, 0, len(request.GetItems()))
	for _, source := range request.GetItems() {
		itemID, err := uuid.Parse(source.GetMenuItemId())
		if err != nil {
			return model.CreateOrder{}, fmt.Errorf("invalid menu item id: %w", err)
		}
		items = append(items, model.CreateOrderItem{
			MenuItemID:      itemID,
			Quantity:        int(source.GetQuantity()),
			ObservedVersion: source.GetObservedVersion(),
		})
	}

	return model.CreateOrder{
		ID:              orderID,
		RestaurantID:    restaurantID,
		CustomerName:    request.GetCustomerName(),
		CustomerPhone:   request.GetCustomerPhone(),
		DeliveryAddress: request.GetDeliveryAddress(),
		Items:           items,
	}, nil
}

func OrderToProto(order model.Order) *partnerpb.RestaurantOrder {
	items := make([]*partnerpb.RestaurantOrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &partnerpb.RestaurantOrderItem{
			MenuItemId:     item.MenuItemID.String(),
			Name:           item.Name,
			UnitPriceMinor: item.UnitPriceMinor,
			Quantity:       int32(item.Quantity),
			LineTotalMinor: item.LineTotalMinor,
		})
	}

	return &partnerpb.RestaurantOrder{
		Id:              order.ID.String(),
		RestaurantId:    order.RestaurantID.String(),
		Status:          orderStatusToProto(order.Status),
		CustomerName:    order.CustomerName,
		CustomerPhone:   order.CustomerPhone,
		DeliveryAddress: order.DeliveryAddress,
		Items:           items,
		TotalMinor:      order.TotalMinor,
		Currency:        order.Currency,
		CreatedAt:       timestamppb.New(order.CreatedAt),
		UpdatedAt:       timestamppb.New(order.UpdatedAt),
	}
}

func orderStatusToProto(status model.OrderStatus) partnerpb.RestaurantOrderStatus {
	switch status {
	case model.OrderStatusAccepted:
		return partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_ACCEPTED
	case model.OrderStatusPreparing:
		return partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_PREPARING
	case model.OrderStatusReady:
		return partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_READY
	case model.OrderStatusDelivering:
		return partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_DELIVERING
	case model.OrderStatusDelivered:
		return partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_DELIVERED
	default:
		return partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_UNSPECIFIED
	}
}
