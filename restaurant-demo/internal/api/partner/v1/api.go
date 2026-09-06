package v1

import (
	"log/slog"

	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

type API struct {
	partnerpb.UnimplementedRestaurantServiceServer
	service RestaurantService
	logger  *slog.Logger
}

func New(service RestaurantService, logger *slog.Logger) *API {
	if logger == nil {
		logger = slog.Default()
	}
	return &API{service: service, logger: logger}
}

var _ partnerpb.RestaurantServiceServer = (*API)(nil)
