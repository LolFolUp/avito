package v1

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

func (a *API) ListRestaurants(ctx context.Context) (openapi.ListRestaurantsRes, error) {
	restaurants, err := a.catalog.ListRestaurants(ctx)
	if err != nil {
		a.logError("list restaurants", err)
		return converter.Problem(err, errorHTTPStatus(err)), nil
	}

	return converter.RestaurantsToOpenAPI(restaurants), nil
}
