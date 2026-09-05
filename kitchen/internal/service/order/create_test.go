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

func TestCreateReturnsMenuChangedBeforeWritingOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	restaurants := mocks.NewMockRestaurantRepository(ctrl)
	menu := mocks.NewMockMenuRepository(ctrl)
	orders := mocks.NewMockOrderRepository(ctrl)
	txManager := mocks.NewMockTxManager(ctrl)
	service := orderservice.New(restaurants, menu, orders, txManager)

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

	txManager.EXPECT().Do(ctx, gomock.Any()).DoAndReturn(
		func(txCtx context.Context, fn func(context.Context) error) error { return fn(txCtx) },
	)
	restaurants.EXPECT().GetRestaurantForUpdate(ctx, restaurantID).Return(model.Restaurant{
		ID: restaurantID, IsActive: true, AcceptingOrders: true, Currency: "RUB",
	}, nil)
	menu.EXPECT().GetMenuItemsForUpdate(ctx, restaurantID, []uuid.UUID{menuItemID}).Return([]model.MenuItem{{
		ID:           menuItemID,
		RestaurantID: restaurantID,
		ExternalID:   "pizza-1",
		Name:         "Margherita",
		PriceMinor:   59000,
		Currency:     "RUB",
		IsActive:     true,
		IsAvailable:  true,
		Version:      2,
	}}, nil)

	_, err := service.Create(ctx, input)
	domainErr, ok := errs.As(err)
	if !ok || domainErr.Code != errs.CodeMenuChanged {
		t.Fatalf("expected MENU_CHANGED, got %v", err)
	}
}
