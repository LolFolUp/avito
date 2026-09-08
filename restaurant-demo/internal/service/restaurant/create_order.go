package restaurant

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (s *Service) CreateOrder(ctx context.Context, input model.CreateOrder) (model.Order, error) {
	input.CustomerName = strings.TrimSpace(input.CustomerName)
	input.CustomerPhone = strings.TrimSpace(input.CustomerPhone)
	input.DeliveryAddress = strings.TrimSpace(input.DeliveryAddress)
	if err := validateCreateOrder(input); err != nil {
		return model.Order{}, err
	}

	var created model.Order
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		itemIDs := make([]uuid.UUID, 0, len(input.Items))
		requestedByID := make(map[uuid.UUID]model.CreateOrderItem, len(input.Items))
		for _, requested := range input.Items {
			itemIDs = append(itemIDs, requested.MenuItemID)
			requestedByID[requested.MenuItemID] = requested
		}

		menuItems, err := s.repository.GetMenuItemsForUpdate(txCtx, itemIDs)
		if err != nil {
			return err
		}
		if len(menuItems) != len(requestedByID) {
			return errs.FailedPrecondition("one or more menu items were removed from the menu")
		}

		orderItems := make([]model.OrderItem, 0, len(menuItems))
		var (
			total         int64
			orderCurrency string
		)
		for _, item := range menuItems {
			requested := requestedByID[item.ID]
			if !item.IsActive ||
				!item.IsAvailable ||
				item.AvailableQuantity < requested.Quantity ||
				item.Version != requested.ObservedVersion {
				return errs.FailedPrecondition("one or more menu items have changed or are unavailable")
			}

			if orderCurrency == "" {
				orderCurrency = item.Currency
			} else if item.Currency != orderCurrency {
				return errs.FailedPrecondition("order items have different currencies")
			}
			if item.PriceMinor > math.MaxInt64/int64(requested.Quantity) {
				return errs.InvalidArgument("order total is too large")
			}

			lineTotal := item.PriceMinor * int64(requested.Quantity)
			if total > math.MaxInt64-lineTotal {
				return errs.InvalidArgument("order total is too large")
			}
			total += lineTotal
			orderItems = append(orderItems, model.OrderItem{
				MenuItemID:     item.ID,
				Name:           item.Name,
				Price:          item.PriceMinor,
				Quantity:       requested.Quantity,
				LineTotalPrice: lineTotal,
			})
		}

		for _, requested := range input.Items {
			if err := s.repository.ReserveMenuItem(
				txCtx,
				requested.MenuItemID,
				requested.Quantity,
			); err != nil {
				return err
			}
		}

		now := time.Now().UTC()
		created = model.Order{
			ID:              input.ID,
			Status:          model.OrderStatusAccepted,
			CustomerName:    input.CustomerName,
			CustomerPhone:   input.CustomerPhone,
			DeliveryAddress: input.DeliveryAddress,
			Items:           orderItems,
			TotalPrice:      total,
			Currency:        orderCurrency,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if err := s.repository.CreateOrder(txCtx, created); err != nil {
			return err
		}
		return s.repository.CreateOrderStatusHistory(txCtx, model.OrderStatusHistory{
			OrderID:  created.ID,
			ToStatus: model.OrderStatusAccepted,
		})
	})
	if err != nil {
		return model.Order{}, err
	}
	return created, nil
}

func validateCreateOrder(input model.CreateOrder) error {
	switch {
	case input.ID == uuid.Nil:
		return errs.InvalidArgument("ID заказа обязателен")
	case input.CustomerName == "" || len(input.CustomerName) > 200:
		return errs.InvalidArgument("Имя недействительно")
	case len(input.CustomerPhone) < 5 || len(input.CustomerPhone) > 32:
		return errs.InvalidArgument("Номер телефона недействителен")
	case input.DeliveryAddress == "" || len(input.DeliveryAddress) > 1000:
		return errs.InvalidArgument("Адрес доставки недействителен.")
	case len(input.Items) == 0 || len(input.Items) > 100:
		return errs.InvalidArgument("Недопустимое количество позиций в заказе")
	}

	seen := make(map[uuid.UUID]struct{}, len(input.Items))
	for _, item := range input.Items {
		if item.MenuItemID == uuid.Nil || item.Quantity < 1 || item.Quantity > 99 || item.ObservedVersion < 1 {
			return errs.InvalidArgument("Позиция в заказе недействителен")
		}
		if _, exists := seen[item.MenuItemID]; exists {
			return errs.InvalidArgument("order contains duplicate menu items")
		}
		seen[item.MenuItemID] = struct{}{}
	}
	return nil
}
