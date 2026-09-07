package restaurant

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (s *Service) GetMenu(ctx context.Context) ([]model.MenuItem, error) {
	menu, err := s.repository.ListMenu(ctx)
	if err != nil {
		return []model.MenuItem{}, err
	}
	return menu, nil

}
