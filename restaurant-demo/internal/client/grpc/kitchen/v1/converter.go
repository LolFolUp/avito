package v1

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func updateOrderStatusToProto(input model.UpdateOrderStatus) *partnerpb.UpdateOrderStatusRequest {
	return &partnerpb.UpdateOrderStatusRequest{
		RestaurantId:   input.RestaurantID.String(),
		OrderId:        input.OrderID.String(),
		ExpectedStatus: orderStatusToProto(input.ExpectedStatus),
		NewStatus:      orderStatusToProto(input.NewStatus),
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
