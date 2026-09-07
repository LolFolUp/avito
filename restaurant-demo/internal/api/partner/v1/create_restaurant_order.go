package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/api/converter"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) CreateRestaurantOrder(
	ctx context.Context,
	request *partnerpb.CreateRestaurantOrderRequest,
) (*partnerpb.CreateRestaurantOrderResponse, error) {
	if err := a.validateRestaurantID(request.GetRestaurantId()); err != nil {
		return nil, grpcError(err)
	}

	input, err := converter.CreateOrderFromProto(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	order, err := a.service.CreateOrder(ctx, input)
	if err != nil {
		a.logError("create restaurant order", err)
		return nil, grpcError(err)
	}
	return &partnerpb.CreateRestaurantOrderResponse{
		Order: converter.OrderToProto(a.restaurantID, order),
	}, nil
}
