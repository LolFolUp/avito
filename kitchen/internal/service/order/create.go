package order

import (
	"context"
	"math"
	"strings"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func (s *Service) Create(ctx context.Context, input model.CreateOrder) (model.Order, error) {
	input.CustomerName = strings.TrimSpace(input.CustomerName)
	input.CustomerPhone = strings.TrimSpace(input.CustomerPhone)
	input.DeliveryAddress = strings.TrimSpace(input.DeliveryAddress)

	if err := validateCreateOrder(input); err != nil {
		return model.Order{}, err
	}

	var created model.Order
	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		restaurant, err := s.restaurants.GetRestaurantForUpdate(txCtx, input.RestaurantID)
		if err != nil {
			return err
		}
		if !restaurant.IsActive || !restaurant.AcceptingOrders {
			return errs.FailedPrecondition(
				errs.CodeRestaurantNotAccepting,
				"restaurant is not accepting orders",
				nil,
			)
		}

		itemIDs := make([]uuid.UUID, 0, len(input.Items))
		requestedByID := make(map[uuid.UUID]model.CreateOrderItem, len(input.Items))
		for _, requested := range input.Items {
			itemIDs = append(itemIDs, requested.MenuItemID)
			requestedByID[requested.MenuItemID] = requested
		}

		menuItems, err := s.menu.GetMenuItemsForUpdate(txCtx, input.RestaurantID, itemIDs)
		if err != nil {
			return err
		}
		if len(menuItems) != len(requestedByID) {
			return errs.NotFound("menu item")
		}

		conflicts := make([]model.MenuItemConflict, 0)
		orderID := uuid.New()
		orderItems := make([]model.OrderItem, 0, len(menuItems))
		var total int64

		for _, item := range menuItems {
			requested := requestedByID[item.ID]
			if !item.IsActive || !item.IsAvailable || item.Version != requested.ObservedVersion {
				conflicts = append(conflicts, model.MenuItemConflict{
					MenuItemID:        item.ID,
					RequestedVersion:  requested.ObservedVersion,
					CurrentVersion:    item.Version,
					IsAvailable:       item.IsActive && item.IsAvailable,
					CurrentPriceMinor: item.PriceMinor,
				})
				continue
			}
			if item.Currency != restaurant.Currency {
				return errs.Internal("validate menu currency", nil)
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
				ID:             uuid.New(),
				OrderID:        orderID,
				MenuItemID:     item.ID,
				ExternalItemID: item.ExternalID,
				Name:           item.Name,
				UnitPriceMinor: item.PriceMinor,
				Quantity:       requested.Quantity,
				LineTotalMinor: lineTotal,
			})
		}

		if len(conflicts) > 0 {
			return errs.FailedPrecondition(
				errs.CodeMenuChanged,
				"one or more menu items have changed",
				conflicts,
			)
		}

		order := model.Order{
			ID:              orderID,
			RestaurantID:    restaurant.ID,
			Status:          model.OrderStatusPendingConfirmation,
			CustomerName:    input.CustomerName,
			CustomerPhone:   input.CustomerPhone,
			DeliveryAddress: input.DeliveryAddress,
			Items:           orderItems,
			TotalMinor:      total,
			Currency:        restaurant.Currency,
		}

		if err := s.orders.CreateOrder(txCtx, order); err != nil {
			return err
		}
		if err := s.orders.CreateOrderItems(txCtx, orderItems); err != nil {
			return err
		}
		if err := s.orders.CreateOrderStatusHistory(txCtx, model.OrderStatusHistory{
			OrderID:  order.ID,
			ToStatus: model.OrderStatusPendingConfirmation,
		}); err != nil {
			return err
		}

		created, err = s.orders.GetOrder(txCtx, order.ID)
		return err
	})
	if err != nil {
		return model.Order{}, err
	}

	return created, nil
}

func validateCreateOrder(input model.CreateOrder) error {
	switch {
	case input.RestaurantID == uuid.Nil:
		return errs.InvalidArgument("restaurant id is required")
	case input.CustomerName == "" || len(input.CustomerName) > 200:
		return errs.InvalidArgument("customer name is invalid")
	case len(input.CustomerPhone) < 5 || len(input.CustomerPhone) > 32:
		return errs.InvalidArgument("customer phone is invalid")
	case input.DeliveryAddress == "" || len(input.DeliveryAddress) > 1000:
		return errs.InvalidArgument("delivery address is invalid")
	case len(input.Items) == 0 || len(input.Items) > 100:
		return errs.InvalidArgument("order items count is invalid")
	}

	seen := make(map[uuid.UUID]struct{}, len(input.Items))
	for _, item := range input.Items {
		if item.MenuItemID == uuid.Nil || item.Quantity < 1 || item.Quantity > 99 || item.ObservedVersion < 1 {
			return errs.InvalidArgument("order item is invalid")
		}
		if _, exists := seen[item.MenuItemID]; exists {
			return errs.InvalidArgument("order contains duplicate menu items")
		}
		seen[item.MenuItemID] = struct{}{}
	}

	return nil
}
