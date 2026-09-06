package record

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID              uuid.UUID
	RestaurantID    uuid.UUID
	Status          string
	CustomerName    string
	CustomerPhone   string
	DeliveryAddress string
	Items           []OrderItem
	TotalMinor      int64
	Currency        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	MenuItemID     uuid.UUID
	Name           string
	UnitPriceMinor int64
	Quantity       int
	LineTotalMinor int64
}
