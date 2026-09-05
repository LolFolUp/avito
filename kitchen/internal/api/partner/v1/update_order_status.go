package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) UpdateOrderStatus(
	ctx context.Context,
	request *partnerpb.UpdateOrderStatusRequest,
) (*partnerpb.UpdateOrderStatusResponse, error) {
	if request == nil {
		return nil, grpcError(errs.InvalidArgument("request is required"))
	}

	restaurantID, err := uuid.Parse(request.GetRestaurantId())
	if err != nil {
		return nil, grpcError(errs.InvalidArgument("restaurant id is invalid"))
	}
	orderID, err := uuid.Parse(request.GetOrderId())
	if err != nil {
		return nil, grpcError(errs.InvalidArgument("order id is invalid"))
	}

	order, err := a.orders.UpdateStatus(
		ctx,
		restaurantID,
		orderID,
		converter.OrderStatusFromProto(request.GetExpectedStatus()),
		converter.OrderStatusFromProto(request.GetNewStatus()),
		request.GetReason(),
	)
	if err != nil {
		a.logError("update order status", err)
		return nil, grpcError(err)
	}

	return &partnerpb.UpdateOrderStatusResponse{Order: converter.OrderToProto(order)}, nil
}
