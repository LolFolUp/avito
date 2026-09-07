package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus описывает этап жизненного цикла заказа.
type OrderStatus string

const (
	// OrderStatusAccepted означает, что ресторан принял заказ.
	OrderStatusAccepted OrderStatus = "ACCEPTED"

	// OrderStatusPreparing означает, что ресторан готовит заказ.
	OrderStatusPreparing OrderStatus = "PREPARING"

	// OrderStatusReady означает, что заказ приготовлен и готов к передаче.
	OrderStatusReady OrderStatus = "READY"

	// OrderStatusDelivering означает, что заказ находится в доставке.
	OrderStatusDelivering OrderStatus = "DELIVERING"

	// OrderStatusDelivered означает, что заказ доставлен.
	// Это терминальный статус.
	OrderStatusDelivered OrderStatus = "DELIVERED"
)

// Order описывает заказ, полученный рестораном из Kitchen.
// Один Order содержит одну или несколько позиций OrderItem.
type Order struct {
	// ID — уникальный идентификатор заказа в Kitchen.
	ID uuid.UUID

	// Status — текущее состояние заказа.
	Status OrderStatus

	// CustomerName — имя клиента, указанное при оформлении.
	CustomerName string

	// CustomerPhone — телефон клиента для связи по заказу.
	CustomerPhone string

	// DeliveryAddress — адрес доставки заказа.
	DeliveryAddress string

	// Items — позиции заказа. Это связь «один заказ — много позиций».
	Items []OrderItem

	// TotalPrice — полная стоимость заказа в минимальных единицах валюты.
	TotalPrice int64

	// Currency — трёхбуквенный код валюты, например RUB.
	Currency string

	// CreatedAt — момент создания заказа рестораном.
	CreatedAt time.Time

	// UpdatedAt — момент последнего изменения заказа.
	UpdatedAt time.Time
}

// OrderItem описывает одну позицию заказа.
// Название и цены являются snapshots на момент оформления заказа.
type OrderItem struct {
	// MenuItemID — идентификатор позиции меню в restaurant-demo.
	MenuItemID uuid.UUID

	// Name — snapshot названия блюда на момент создания заказа.
	Name string

	// Price — snapshot цены одной единицы блюда
	// в минимальных единицах валюты.
	Price int64

	// Quantity — количество заказанных единиц блюда.
	Quantity int

	// LineTotalPrice — стоимость позиции: Price * Quantity.
	LineTotalPrice int64
}

type CreateOrderItem struct {
	MenuItemID      uuid.UUID
	Quantity        int
	ObservedVersion int64
}

type CreateOrder struct {
	ID              uuid.UUID
	CustomerName    string
	CustomerPhone   string
	DeliveryAddress string
	Items           []CreateOrderItem
}

// UpdateOrderStatus задаёт ожидаемый и новый статусы заказа.
type UpdateOrderStatus struct {
	// RestaurantID — идентификатор этого экземпляра restaurant-demo в Kitchen.
	RestaurantID uuid.UUID
	// OrderID — идентификатор изменяемого заказа.
	OrderID uuid.UUID
	// ExpectedStatus — статус, в котором заказ должен находиться до изменения.
	ExpectedStatus OrderStatus
	// NewStatus — статус, в который требуется перевести заказ.
	NewStatus OrderStatus
}

// OrderStatusChange содержит результат смены статуса ресторанного заказа.
type OrderStatusChange struct {
	// OrderID — идентификатор изменённого заказа.
	OrderID uuid.UUID
	// Status — новый статус заказа.
	Status OrderStatus
	// UpdatedAt — момент изменения заказа в restaurant-demo.
	UpdatedAt time.Time
}

// OrderStatusHistory описывает один переход статуса заказа.
type OrderStatusHistory struct {
	OrderID    uuid.UUID
	FromStatus *OrderStatus
	ToStatus   OrderStatus
	CreatedAt  time.Time
}
