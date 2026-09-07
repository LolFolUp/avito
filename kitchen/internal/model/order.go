package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus описывает текущее состояние заказа в Kitchen.
type OrderStatus string

const (
	// OrderStatusPendingConfirmation означает, что заказ ожидает решения ресторана.
	OrderStatusPendingConfirmation OrderStatus = "PENDING_CONFIRMATION"
	// OrderStatusAccepted означает, что ресторан принял заказ.
	OrderStatusAccepted OrderStatus = "ACCEPTED"
	// OrderStatusPreparing означает, что ресторан готовит заказ.
	OrderStatusPreparing OrderStatus = "PREPARING"
	// OrderStatusReady означает, что заказ приготовлен и готов к передаче.
	OrderStatusReady OrderStatus = "READY"
	// OrderStatusDelivering означает, что заказ передан в доставку.
	OrderStatusDelivering OrderStatus = "DELIVERING"
	// OrderStatusDelivered означает, что заказ доставлен.
	OrderStatusDelivered OrderStatus = "DELIVERED"
	// OrderStatusRejected означает, что ресторан отклонил заказ.
	OrderStatusRejected OrderStatus = "REJECTED"
)

// Order содержит сохранённое в Kitchen публичное состояние заказа.
type Order struct {
	// ID — уникальный идентификатор заказа, общий для Kitchen и restaurant-demo.
	ID uuid.UUID
	// RestaurantID — идентификатор ресторана из каталога Kitchen.
	RestaurantID uuid.UUID
	// Status — текущее состояние заказа.
	Status OrderStatus
	// CustomerName — имя клиента, указанное при оформлении.
	CustomerName string
	// CustomerPhone — телефон клиента для связи по заказу.
	CustomerPhone string
	// DeliveryAddress — адрес доставки заказа.
	DeliveryAddress string
	// Items — позиции заказа со snapshots названий и цен.
	Items []OrderItem
	// TotalPrice — полная стоимость заказа в минимальных единицах валюты.
	TotalPrice int64
	// Currency — трёхбуквенный код валюты, например RUB.
	Currency string
	// RejectionReason — причина отказа ресторана; пустая строка означает отсутствие причины.
	RejectionReason string
	// CreatedAt — момент создания заказа.
	CreatedAt time.Time
	// UpdatedAt — момент последнего изменения заказа.
	UpdatedAt time.Time
}

// OrderItem описывает одну позицию заказа и её неизменяемый snapshot.
type OrderItem struct {
	// OrderID — идентификатор заказа, которому принадлежит позиция.
	OrderID uuid.UUID
	// MenuItemID — идентификатор блюда в базе restaurant-demo.
	MenuItemID uuid.UUID
	// Name — название блюда на момент создания заказа.
	Name string
	// UnitPriceMinor — цена одной единицы блюда в минимальных единицах валюты.
	UnitPriceMinor int64
	// Quantity — заказанное количество блюда.
	Quantity int
	// LineTotalPrice — стоимость позиции: UnitPriceMinor * Quantity.
	LineTotalPrice int64
}

// CreateOrderItem описывает выбранную клиентом позицию нового заказа.
type CreateOrderItem struct {
	// MenuItemID — идентификатор блюда, полученный из live-меню ресторана.
	MenuItemID uuid.UUID
	// Quantity — количество единиц блюда.
	Quantity int
	// ObservedVersion — версия блюда, которую клиент увидел при чтении меню.
	ObservedVersion int64
}

// CreateOrder содержит данные, необходимые для оформления нового заказа.
type CreateOrder struct {
	// RestaurantID — ресторан, в котором клиент оформляет заказ.
	RestaurantID uuid.UUID
	// CustomerName — имя клиента.
	CustomerName string
	// CustomerPhone — телефон клиента.
	CustomerPhone string
	// DeliveryAddress — адрес доставки.
	DeliveryAddress string
	// Items — выбранные клиентом позиции меню.
	Items []CreateOrderItem
}

// OrderStatusHistory описывает одну запись истории изменения статуса заказа.
type OrderStatusHistory struct {
	// ID — внутренний последовательный идентификатор записи истории.
	ID int64
	// OrderID — идентификатор заказа.
	OrderID uuid.UUID
	// Status — статус, установленный этой записью.
	Status OrderStatus
	// Reason — причина изменения статуса; используется при отказе.
	Reason string
	// CreatedAt — момент изменения статуса.
	CreatedAt time.Time
}
