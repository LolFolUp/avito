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
		return model.CreateOrder{}, fmt.Errorf("невалидный ID заказа: %w", err)
	}

	items := make([]model.CreateOrderItem, 0, len(request.GetItems()))
	for _, source := range request.GetItems() {
		itemID, err := uuid.Parse(source.GetMenuItemId())
		if err != nil {
			return model.CreateOrder{}, fmt.Errorf("невалидный ID позиции меню: %w", err)
		}
		items = append(items, model.CreateOrderItem{
			MenuItemID:      itemID,
			Quantity:        int(source.GetQuantity()),
			ObservedVersion: source.GetObservedVersion(),
		})
	}

	return model.CreateOrder{
		ID:              orderID,
		CustomerName:    request.GetCustomerName(),
		CustomerPhone:   request.GetCustomerPhone(),
		DeliveryAddress: request.GetDeliveryAddress(),
		Items:           items,
	}, nil
}

func OrderToProto(restaurantID uuid.UUID, order model.Order) *partnerpb.RestaurantOrder {
	items := make([]*partnerpb.RestaurantOrderItem, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &partnerpb.RestaurantOrderItem{
			MenuItemId:     item.MenuItemID.String(),
			Name:           item.Name,
			UnitPriceMinor: item.Price,
			Quantity:       int32(item.Quantity),
			LineTotalMinor: item.LineTotalPrice,
		})
	}

	return &partnerpb.RestaurantOrder{
		Id:              order.ID.String(),
		RestaurantId:    restaurantID.String(),
		Status:          orderStatusToProto(order.Status),
		CustomerName:    order.CustomerName,
		CustomerPhone:   order.CustomerPhone,
		DeliveryAddress: order.DeliveryAddress,
		Items:           items,
		TotalMinor:      order.TotalPrice,
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

func UpdateOrderStatusFromProto(
	request *partnerpb.UpdateRestaurantOrderStatusRequest,
) (model.UpdateOrderStatus, error) {
	restaurantID, err := uuid.Parse(request.GetRestaurantId())
	if err != nil {
		return model.UpdateOrderStatus{}, fmt.Errorf("restaurant id must be a UUID: %w", err)
	}
	orderID, err := uuid.Parse(request.GetOrderId())
	if err != nil {
		return model.UpdateOrderStatus{}, fmt.Errorf("order id must be a UUID: %w", err)
	}
	expectedStatus, err := orderStatusFromProto(request.GetExpectedStatus())
	if err != nil {
		return model.UpdateOrderStatus{}, fmt.Errorf("expected status: %w", err)
	}
	newStatus, err := orderStatusFromProto(request.GetNewStatus())
	if err != nil {
		return model.UpdateOrderStatus{}, fmt.Errorf("new status: %w", err)
	}
	return model.UpdateOrderStatus{
		RestaurantID:   restaurantID,
		OrderID:        orderID,
		ExpectedStatus: expectedStatus,
		NewStatus:      newStatus,
	}, nil
}

func OrderStatusChangeToProto(change model.OrderStatusChange) *partnerpb.UpdateRestaurantOrderStatusResponse {
	return &partnerpb.UpdateRestaurantOrderStatusResponse{
		OrderId:   change.OrderID.String(),
		Status:    orderStatusToProto(change.Status),
		UpdatedAt: timestamppb.New(change.UpdatedAt),
	}
}

func orderStatusFromProto(status partnerpb.RestaurantOrderStatus) (model.OrderStatus, error) {
	switch status {
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_ACCEPTED:
		return model.OrderStatusAccepted, nil
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_PREPARING:
		return model.OrderStatusPreparing, nil
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_READY:
		return model.OrderStatusReady, nil
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_DELIVERING:
		return model.OrderStatusDelivering, nil
	case partnerpb.RestaurantOrderStatus_RESTAURANT_ORDER_STATUS_DELIVERED:
		return model.OrderStatusDelivered, nil
	default:
		return "", fmt.Errorf("status must be specified")
	}
}
