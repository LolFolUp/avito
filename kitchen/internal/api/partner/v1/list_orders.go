package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) ListOrders(ctx context.Context, request *partnerpb.ListOrdersRequest) (*partnerpb.ListOrdersResponse, error) {
	if request == nil {
		return nil, grpcError(errs.InvalidArgument("request is required"))
	}

	restaurantID, err := uuid.Parse(request.GetRestaurantId())
	if err != nil {
		return nil, grpcError(errs.InvalidArgument("restaurant id is invalid"))
	}

	page, err := a.orders.List(
		ctx,
		restaurantID,
		converter.OrderStatusFromProto(request.GetStatus()),
		int(request.GetPageSize()),
		request.GetPageToken(),
	)
	if err != nil {
		a.logError("list orders", err)
		return nil, grpcError(err)
	}

	return &partnerpb.ListOrdersResponse{
		Orders:        converter.OrdersToProto(page.Orders),
		NextPageToken: page.NextPageToken,
	}, nil
}
