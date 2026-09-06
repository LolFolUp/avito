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

// OrderItem описывает одну позицию заказа.
// Название и цены являются snapshots на момент оформления заказа.
type OrderItem struct {
	MenuItemID     uuid.UUID
	Name           string
	UnitPriceMinor int64
	Quantity       int
	LineTotalMinor int64
}

type CreateOrderItem struct {
	MenuItemID      uuid.UUID
	Quantity        int
	ObservedVersion int64
}
