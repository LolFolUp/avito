package restaurant

import (
	"log/slog"

	"github.com/google/uuid"
)

type Service struct {
	restaurantID uuid.UUID
	currency     string
	repository   RestaurantRepository
	txManager    TxManager
	logger       *slog.Logger
}

func New(
	restaurantID uuid.UUID,
	currency string,
	repository RestaurantRepository,
	txManager TxManager,
	logger *slog.Logger,
) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		restaurantID: restaurantID,
		currency:     currency,
		repository:   repository,
		txManager:    txManager,
		logger:       logger,
	}
}
