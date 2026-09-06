package model

import "github.com/google/uuid"

type MenuItem struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Name         string
	Description  string
	PriceMinor   int64
	Currency     string
	IsAvailable  bool
	Version      int64
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
