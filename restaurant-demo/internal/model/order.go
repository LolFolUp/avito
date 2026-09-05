package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPendingConfirmation OrderStatus = "pending_confirmation"
	OrderStatusAccepted            OrderStatus = "accepted"
	OrderStatusPreparing           OrderStatus = "preparing"
	OrderStatusReady               OrderStatus = "ready"
	OrderStatusDelivering          OrderStatus = "delivering"
	OrderStatusDelivered           OrderStatus = "delivered"
	OrderStatusRejected            OrderStatus = "rejected"
)

type Order struct {
	ID              uuid.UUID
	RestaurantID    uuid.UUID
	Status          OrderStatus
	CustomerName    string
	CustomerPhone   string
	DeliveryAddress string
	Items           []OrderItem
	TotalMinor      int64
	Currency        string
	RejectionReason string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderItem struct {
	MenuItemID     uuid.UUID
	ExternalItemID string
	Name           string
	UnitPriceMinor int64
	Quantity       int
	LineTotalMinor int64
}

type OrderPage struct {
	Orders        []Order
	NextPageToken string
}
