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
	TotalPrice      int64
	Currency        string
	RejectionReason pgtype.Text // Тут есть еще варианты *string(как я обычно делал) и sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	OrderID        uuid.UUID
	MenuItemID     uuid.UUID
	Name           string
	UnitPriceMinor int64
	Quantity       int
	LineTotalPrice int64
}

type OrderStatusHistory struct {
	ID        int64
	OrderID   uuid.UUID
	Status    string
	Reason    pgtype.Text
	CreatedAt time.Time
}
