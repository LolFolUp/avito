package v1

import (
	"log/slog"

	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

type API struct {
	partnerpb.UnimplementedKitchenPartnerServiceServer
	orders OrderService
	logger *slog.Logger
}

func New(orders OrderService, logger *slog.Logger) *API {
	if logger == nil {
		logger = slog.Default()
	}
	return &API{orders: orders, logger: logger}
}

var _ partnerpb.KitchenPartnerServiceServer = (*API)(nil)
