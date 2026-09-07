package model

import "github.com/google/uuid"

// MenuItem описывает позицию live-меню, полученную от restaurant-demo.
type MenuItem struct {
	// ID — идентификатор позиции в базе restaurant-demo.
	ID uuid.UUID

	// RestaurantID — идентификатор ресторана, которому принадлежит позиция.
	RestaurantID uuid.UUID

	// Name — отображаемое название блюда.
	Name string

	// Description — необязательное описание блюда.
	Description string

	// PriceMinor — цена в минимальных единицах валюты.
	PriceMinor int64

	// Currency — трёхбуквенный код валюты, например RUB.
	Currency string

	// IsAvailable показывает, можно ли заказать позицию прямо сейчас.
	IsAvailable bool

	// Version изменяется при значимом изменении позиции и защищает создание заказа от устаревшего меню.
	Version int64
}

// RestaurantMenu объединяет сведения о ресторане и полученное от него live-меню.
type RestaurantMenu struct {
	// Restaurant — ресторан из каталога Kitchen.
	Restaurant Restaurant
	// Items — актуальные позиции, возвращённые restaurant-demo.
	Items []MenuItem
}

// MenuItemConflict содержит сведения о позиции, изменившейся после просмотра меню.
type MenuItemConflict struct {
	// MenuItemID — идентификатор конфликтующей позиции.
	MenuItemID uuid.UUID

	// RequestedVersion — версия, присланная клиентом.
	RequestedVersion int64

	// CurrentVersion — актуальная версия позиции в restaurant-demo.
	CurrentVersion int64

	// IsAvailable — актуальная доступность позиции.
	IsAvailable bool

	// CurrentPriceMinor — актуальная цена позиции в минимальных единицах валюты.
	CurrentPriceMinor int64
}
