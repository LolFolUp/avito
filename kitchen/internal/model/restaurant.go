package model

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID              uuid.UUID
	Name            string
	IsActive        bool
	AcceptingOrders bool
	Currency        string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
