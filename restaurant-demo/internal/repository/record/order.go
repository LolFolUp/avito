package record

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID
	Status          string
	CustomerName    string
	CustomerPhone   string
	DeliveryAddress string
	TotalPrice      int64
	Currency        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// OrderItem описывает одну позицию заказа.
// Название и цены являются snapshots на момент оформления заказа.
type OrderItem struct {
	OrderID        uuid.UUID
	MenuItemID     uuid.UUID
	Name           string
	Price          int64
	Quantity       int
	LineTotalPrice int64
}

type CreateOrderItem struct {
	MenuItemID      uuid.UUID
	Quantity        int
	ObservedVersion int64
}
