package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/api/converter"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) UpdateRestaurantOrderStatus(
	ctx context.Context,
	request *partnerpb.UpdateRestaurantOrderStatusRequest,
) (*partnerpb.UpdateRestaurantOrderStatusResponse, error) {
	if err := a.validateRestaurantID(request.GetRestaurantId()); err != nil {
		return nil, grpcError(err)
	}

	input, err := converter.UpdateOrderStatusFromProto(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	change, err := a.service.UpdateOrderStatus(ctx, input)
	if err != nil {
		a.logError("update restaurant order status", err)
		return nil, grpcError(err)
	}
	return converter.OrderStatusChangeToProto(change), nil
}
