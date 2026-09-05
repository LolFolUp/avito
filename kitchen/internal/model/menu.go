package model

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

type MenuSyncResult struct {
	CreatedCount     int
	UpdatedCount     int
	DeactivatedCount int
}

type RestaurantMenu struct {
	Restaurant Restaurant
	Items      []MenuItem
}

type MenuItemConflict struct {
	MenuItemID        uuid.UUID
	RequestedVersion  int64
	CurrentVersion    int64
	IsAvailable       bool
	CurrentPriceMinor int64
}
