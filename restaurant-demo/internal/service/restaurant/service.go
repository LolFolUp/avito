package restaurant

import (
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	restaurantID uuid.UUID
	partner      PartnerClient
	menu         MenuRepository
	logger       *slog.Logger
}

func New(restaurantID uuid.UUID, partner PartnerClient, menu MenuRepository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		restaurantID: restaurantID,
		partner:      partner,
		menu:         menu,
		logger:       logger,
	}
}
