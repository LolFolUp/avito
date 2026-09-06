package order

import (
	"context"
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

	restaurant, err := s.restaurants.GetRestaurant(ctx, input.RestaurantID)
	if err != nil {
		return model.Order{}, err
	}
	if !restaurant.IsActive || !restaurant.AcceptingOrders {
		return model.Order{}, errs.FailedPrecondition(
			errs.CodeRestaurantNotAccepting,
			"restaurant is not accepting orders",
			nil,
		)
	}

	orderID := uuid.New()
	accepted, err := s.partner.CreateOrder(ctx, orderID, input)
	if err != nil {
		return model.Order{}, err
	}
	if accepted.ID != orderID ||
		accepted.RestaurantID != input.RestaurantID ||
		accepted.Status != model.OrderStatusAccepted ||
		accepted.Currency != restaurant.Currency {
		return model.Order{}, errs.Internal("validate restaurant order response", nil)
	}

	var created model.Order
	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := s.orders.CreateOrder(txCtx, accepted); err != nil {
			return err
		}
		if err := s.orders.CreateOrderItems(txCtx, accepted.Items); err != nil {
			return err
		}
		if err := s.orders.CreateOrderStatusHistory(txCtx, model.OrderStatusHistory{
			OrderID: accepted.ID,
			Status:  model.OrderStatusAccepted,
		}); err != nil {
			return err
		}

		var err error
		created, err = s.orders.GetOrder(txCtx, accepted.ID)
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
