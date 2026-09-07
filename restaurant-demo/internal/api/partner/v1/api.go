package v1

import (
	"log/slog"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

type API struct {
	partnerpb.UnimplementedRestaurantServiceServer
	restaurantID uuid.UUID
	service      RestaurantService
	logger       *slog.Logger
}

func New(restaurantID uuid.UUID, service RestaurantService, logger *slog.Logger) *API {
	if logger == nil {
		logger = slog.Default()
	}
	return &API{restaurantID: restaurantID, service: service, logger: logger}
}

func (a *API) validateRestaurantID(value string) error {
	restaurantID, err := uuid.Parse(value)
	if err != nil {
		return errs.InvalidArgument("restaurant id must be a UUID")
	}
	if restaurantID != a.restaurantID {
		return errs.NotFound("restaurant not found")
	}
	return nil
}

var _ partnerpb.RestaurantServiceServer = (*API)(nil)
