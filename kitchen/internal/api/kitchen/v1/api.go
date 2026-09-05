package v1

import (
	"log/slog"

	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

type API struct {
	catalog   CatalogService
	orders    OrderService
	readiness ReadinessChecker
	logger    *slog.Logger
}

var _ openapi.Handler = (*API)(nil)

func New(catalog CatalogService, orders OrderService, readiness ReadinessChecker, logger *slog.Logger) *API {
	if logger == nil {
		logger = slog.Default()
	}
	return &API{
		catalog:   catalog,
		orders:    orders,
		readiness: readiness,
		logger:    logger,
	}
}
