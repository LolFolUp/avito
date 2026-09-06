package v1

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/api/converter"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) GetMenu(ctx context.Context, request *partnerpb.GetMenuRequest) (*partnerpb.GetMenuResponse, error) {
	restaurantID, err := uuid.Parse(request.GetRestaurantId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "restaurant id must be a UUID")
	}

	items, currency, err := a.service.GetMenu(ctx, restaurantID)
	if err != nil {
		a.logError("get restaurant menu", err)
		return nil, grpcError(err)
	}
	return &partnerpb.GetMenuResponse{
		RestaurantId: restaurantID.String(),
		Currency:     currency,
		Items:        converter.MenuItemsToProto(items),
	}, nil
}
