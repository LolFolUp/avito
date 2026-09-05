package record

import (
	"time"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	ExternalID   string
	Name         string
	Description  string
	PriceMinor   int64
	Currency     string
	IsActive     bool
	IsAvailable  bool
	Version      int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
