package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func UpdateOrderStatusFromProto(
	request *partnerpb.UpdateOrderStatusRequest,
) (uuid.UUID, uuid.UUID, model.OrderStatus, model.OrderStatus, error) {
	restaurantID, err := uuid.Parse(request.GetRestaurantId())
	if err != nil {
		return uuid.Nil, uuid.Nil, "", "", fmt.Errorf("restaurant id must be a UUID: %w", err)
	}
	orderID, err := uuid.Parse(request.GetOrderId())
	if err != nil {
		return uuid.Nil, uuid.Nil, "", "", fmt.Errorf("order id must be a UUID: %w", err)
	}
	expectedStatus, err := OrderStatusFromProto(request.GetExpectedStatus())
	if err != nil {
		return uuid.Nil, uuid.Nil, "", "", fmt.Errorf("expected status: %w", err)
	}
	newStatus, err := OrderStatusFromProto(request.GetNewStatus())
	if err != nil {
		return uuid.Nil, uuid.Nil, "", "", fmt.Errorf("new status: %w", err)
	}
	return restaurantID, orderID, expectedStatus, newStatus, nil
}

func OrderStatusFromProto(status partnerpb.RestaurantOrderStatus) (model.OrderStatus, error) {
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

func OrderStatusToProto(status model.OrderStatus) partnerpb.RestaurantOrderStatus {
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
