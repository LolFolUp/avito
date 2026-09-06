package v1

import (
	"log/slog"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

func (a *API) logError(operation string, err error) {
	domainErr, ok := errs.As(err)
	if ok && domainErr.Kind != errs.KindInternal {
		return
	}
	a.logger.Error(operation, slog.Any("error", err))
}

func errorHTTPStatus(err error) int {
	domainErr, ok := errs.As(err)
	if !ok {
		return 500
	}

	switch domainErr.Kind {
	case errs.KindInvalidArgument:
		return 400
	case errs.KindNotFound:
		return 404
	case errs.KindAlreadyExists, errs.KindFailedPrecondition:
		return 409
	case errs.KindUnavailable:
		return 503
	default:
		return 500
	}
}

func createOrderError(err error) openapi.CreateOrderRes {
	problem := converter.Error(err, errorHTTPStatus(err))
	switch problem.Status {
	case 400:
		return (*openapi.CreateOrderBadRequest)(problem)
	case 404:
		return (*openapi.CreateOrderNotFound)(problem)
	case 409:
		return (*openapi.CreateOrderConflict)(problem)
	case 503:
		return (*openapi.CreateOrderServiceUnavailable)(problem)
	default:
		return (*openapi.CreateOrderInternalServerError)(problem)
	}
}

func getOrderError(err error) openapi.GetOrderRes {
	problem := converter.Error(err, errorHTTPStatus(err))
	switch problem.Status {
	case 400:
		return (*openapi.GetOrderBadRequest)(problem)
	case 404:
		return (*openapi.GetOrderNotFound)(problem)
	default:
		return (*openapi.GetOrderInternalServerError)(problem)
	}
}

func restaurantMenuError(err error) openapi.GetRestaurantMenuRes {
	problem := converter.Error(err, errorHTTPStatus(err))
	switch problem.Status {
	case 400:
		return (*openapi.GetRestaurantMenuBadRequest)(problem)
	case 404:
		return (*openapi.GetRestaurantMenuNotFound)(problem)
	case 503:
		return (*openapi.GetRestaurantMenuServiceUnavailable)(problem)
	default:
		return (*openapi.GetRestaurantMenuInternalServerError)(problem)
	}
}

func converterError(err error, statusCode int) *openapi.Error {
	return converter.Error(err, statusCode)
}
