package order_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	orderservice "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/order"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/order/mocks"
)

func TestCreateReturnsRestaurantMenuConflictBeforeKitchenTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	restaurants := mocks.NewMockRestaurantRepository(ctrl)
	partner := mocks.NewMockRestaurantClient(ctrl)
	orders := mocks.NewMockOrderRepository(ctrl)
	txManager := mocks.NewMockTxManager(ctrl)
	service := orderservice.New(restaurants, partner, orders, txManager)

	restaurantID := uuid.New()
	menuItemID := uuid.New()
	ctx := context.Background()
	input := model.CreateOrder{
		RestaurantID:    restaurantID,
		CustomerName:    "Alice",
		CustomerPhone:   "+79990000000",
		DeliveryAddress: "Moscow",
		Items: []model.CreateOrderItem{{
			MenuItemID:      menuItemID,
			Quantity:        1,
			ObservedVersion: 1,
		}},
	}

	restaurants.EXPECT().GetRestaurant(ctx, restaurantID).Return(model.Restaurant{
		ID: restaurantID, IsActive: true, AcceptingOrders: true, Currency: "RUB",
	}, nil)
	partner.EXPECT().CreateOrder(ctx, gomock.Any(), input).Return(
		model.Order{},
		errs.FailedPrecondition(errs.CodeMenuChanged, "menu changed", nil),
	)

	_, err := service.Create(ctx, input)
	domainErr, ok := errs.As(err)
	if !ok || domainErr.Code != errs.CodeMenuChanged {
		t.Fatalf("expected MENU_CHANGED, got %v", err)
	}
}
