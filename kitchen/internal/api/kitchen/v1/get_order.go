package v1

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

func (a *API) GetOrder(ctx context.Context, params openapi.GetOrderParams) (openapi.GetOrderRes, error) {
	order, err := a.orders.Get(ctx, params.OrderId)
	if err != nil {
		a.logError("get order", err)
		return getOrderError(err), nil
	}

	return converter.OrderToOpenAPI(order), nil
}
