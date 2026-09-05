package restaurant

import "context"

func (s *Service) SyncMenu(ctx context.Context) error {
	menu, err := s.menu.ListMenu(ctx)
	if err != nil {
		return err
	}

	result, err := s.partner.SyncMenu(ctx, s.restaurantID, menu)
	if err != nil {
		return err
	}

	s.logger.Info(
		"menu synchronized",
		"created", result.CreatedCount,
		"updated", result.UpdatedCount,
		"deactivated", result.DeactivatedCount,
	)
	return nil
}
