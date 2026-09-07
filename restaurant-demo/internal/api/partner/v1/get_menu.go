package v1

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/api/converter"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) GetMenu(ctx context.Context, request *partnerpb.GetMenuRequest) (*partnerpb.GetMenuResponse, error) {
	if err := a.validateRestaurantID(request.GetRestaurantId()); err != nil {
		return nil, grpcError(err)
	}

	items, err := a.service.GetMenu(ctx)
	if err != nil {
		a.logError("get restaurant menu", err)
		return nil, grpcError(err)
	}
	return &partnerpb.GetMenuResponse{
		Items: converter.MenuItemsToProto(items),
	}, nil
}
