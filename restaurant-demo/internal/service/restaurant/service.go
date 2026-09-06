package restaurant

import (
	"log/slog"
)

type Service struct {
	repository RestaurantRepository
	txManager  TxManager
	logger     *slog.Logger
}

func New(repository RestaurantRepository, txManager TxManager, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repository: repository,
		txManager:  txManager,
		logger:     logger,
	}
}
