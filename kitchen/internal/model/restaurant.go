package model

import (
	"time"

	"github.com/google/uuid"
)

// Restaurant описывает заведение, зарегистрированное в каталоге Kitchen.
type Restaurant struct {
	// ID — уникальный идентификатор ресторана.
	ID uuid.UUID
	// Name — отображаемое название ресторана.
	Name string
	// IsActive показывает, доступен ли ресторан в каталоге Kitchen.
	IsActive bool
	// AcceptingOrders показывает, принимает ли ресторан новые заказы.
	AcceptingOrders bool
	// Currency — валюта заказов ресторана в формате ISO 4217, например RUB.
	Currency string
	// CreatedAt — момент регистрации ресторана в Kitchen.
	CreatedAt time.Time
	// UpdatedAt — момент последнего изменения ресторана.
	UpdatedAt time.Time
}
