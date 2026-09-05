package memory

import (
	"context"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (r *Repository) ListMenu(context.Context) ([]model.MenuItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]model.MenuItem, 0, len(r.menu))
	for _, item := range r.menu {
		items = append(items, item)
	}
	return items, nil
}

func (r *Repository) ItemsAvailable(_ context.Context, externalIDs []string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, externalID := range externalIDs {
		item, exists := r.menu[externalID]
		if !exists || !item.IsActive || !item.IsAvailable {
			return false, nil
		}
	}
	return true, nil
}
