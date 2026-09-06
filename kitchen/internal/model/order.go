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
	OrderID        uuid.UUID
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

type CreateOrder struct {
	RestaurantID    uuid.UUID
	CustomerName    string
	CustomerPhone   string
	DeliveryAddress string
	Items           []CreateOrderItem
}

type OrderFilter struct {
	RestaurantID   uuid.UUID
	Status         OrderStatus
	Limit          int
	AfterCreatedAt *time.Time
	AfterID        *uuid.UUID
}

type OrderPage struct {
	Orders        []Order
	NextPageToken string
}

type OrderStatusHistory struct {
	ID        int64
	OrderID   uuid.UUID
	Status    OrderStatus
	Reason    string
	CreatedAt time.Time
}
