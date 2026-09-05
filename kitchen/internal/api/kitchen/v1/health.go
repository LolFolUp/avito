package v1

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

func (a *API) Health(context.Context) (*openapi.HealthStatus, error) {
	return &openapi.HealthStatus{Status: openapi.HealthStatusStatusOk}, nil
}

func (a *API) Readiness(ctx context.Context) (openapi.ReadinessRes, error) {
	if err := a.readiness.Ping(ctx); err != nil {
		a.logError("readiness", err)
		return converterProblem(errs.Internal("check readiness", err), 503), nil
	}

	return &openapi.HealthStatus{Status: openapi.HealthStatusStatusOk}, nil
}
