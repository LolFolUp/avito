package v1

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

func (a *API) GetRestaurantMenu(
	ctx context.Context,
	params openapi.GetRestaurantMenuParams,
) (openapi.GetRestaurantMenuRes, error) {
	menu, err := a.catalog.GetMenu(ctx, params.RestaurantId)
	if err != nil {
		a.logError("get restaurant menu", err)
		return restaurantMenuError(err), nil
	}

	return converter.MenuToOpenAPI(menu), nil
}
