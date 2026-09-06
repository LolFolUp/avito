package model

import "github.com/google/uuid"

// MenuItem описывает позицию меню во внутренней модели restaurant-demo.
type MenuItem struct {
	ID                uuid.UUID
	RestaurantID      uuid.UUID
	Name              string
	Description       string
	PriceMinor        int64
	Currency          string
	IsActive          bool
	IsAvailable       bool
	AvailableQuantity int
	Version           int64
}

type Restaurant struct {
	ID       uuid.UUID
	Name     string
	Currency string
}

type MenuItemConflict struct {
	MenuItemID       uuid.UUID
	RequestedVersion int64
	CurrentVersion   int64
	IsAvailable      bool
}
