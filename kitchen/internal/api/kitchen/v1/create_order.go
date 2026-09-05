package v1

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

func (a *API) CreateOrder(ctx context.Context, request *openapi.CreateOrderRequest) (openapi.CreateOrderRes, error) {
	order, err := a.orders.Create(ctx, converter.CreateOrderToModel(request))
	if err != nil {
		a.logError("create order", err)
		return createOrderError(err), nil
	}

	return converter.OrderToOpenAPI(order), nil
}
