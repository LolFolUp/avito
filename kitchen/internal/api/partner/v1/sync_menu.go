package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) SyncMenu(ctx context.Context, request *partnerpb.SyncMenuRequest) (*partnerpb.SyncMenuResponse, error) {
	if request == nil {
		return nil, grpcError(errs.InvalidArgument("request is required"))
	}

	restaurantID, err := uuid.Parse(request.GetRestaurantId())
	if err != nil {
		return nil, grpcError(errs.InvalidArgument("restaurant id is invalid"))
	}

	result, err := a.catalog.SyncMenu(ctx, restaurantID, converter.MenuItemsFromProto(request.GetItems()))
	if err != nil {
		a.logError("sync menu", err)
		return nil, grpcError(err)
	}

	return &partnerpb.SyncMenuResponse{
		CreatedCount:     int32(result.CreatedCount),
		UpdatedCount:     int32(result.UpdatedCount),
		DeactivatedCount: int32(result.DeactivatedCount),
	}, nil
}
