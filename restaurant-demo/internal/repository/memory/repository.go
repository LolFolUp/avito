package memory

import (
	"sync"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

type Repository struct {
	mu   sync.RWMutex
	menu map[string]model.MenuItem
}

func New(menu []model.MenuItem) *Repository {
	items := make(map[string]model.MenuItem, len(menu))
	for _, item := range menu {
		items[item.ExternalID] = item
	}
	return &Repository{menu: items}
}
