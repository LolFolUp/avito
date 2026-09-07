package restaurant

import (
	"log/slog"
)

type Service struct {
	repository RestaurantRepository
	txManager  TxManager
	kitchen    KitchenClient
	logger     *slog.Logger
}

func New(
	repository RestaurantRepository,
	txManager TxManager,
	kitchen KitchenClient,
	logger *slog.Logger,
) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repository: repository,
		txManager:  txManager,
		kitchen:    kitchen,
		logger:     logger,
	}
}
