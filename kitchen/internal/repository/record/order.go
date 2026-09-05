package record

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID              uuid.UUID
	RestaurantID    uuid.UUID
	Status          string
	CustomerName    string
	CustomerPhone   string
	DeliveryAddress string
	TotalMinor      int64
	Currency        string
	RejectionReason pgtype.Text
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	ID             uuid.UUID
	OrderID        uuid.UUID
	MenuItemID     uuid.UUID
	ExternalItemID string
	NameSnapshot   string
	UnitPriceMinor int64
	Quantity       int
	LineTotalMinor int64
}

type OrderStatusHistory struct {
	ID         int64
	OrderID    uuid.UUID
	FromStatus pgtype.Text
	ToStatus   string
	Reason     pgtype.Text
	CreatedAt  time.Time
}
