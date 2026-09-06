package record

import "github.com/google/uuid"

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
