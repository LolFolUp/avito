package restaurant_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/service/restaurant"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/service/restaurant/mocks"
)

func TestCreateOrderDoesNotReserveWhenObservedVersionIsStale(t *testing.T) {
	ctrl := gomock.NewController(t)
	repository := mocks.NewMockRestaurantRepository(ctrl)
	txManager := mocks.NewMockTxManager(ctrl)

	ctx := context.Background()
	service := restaurant.New(repository, txManager, nil)
	itemID := uuid.New()
	input := model.CreateOrder{
		ID:              uuid.New(),
		CustomerName:    "Alice",
		CustomerPhone:   "+79990000000",
		DeliveryAddress: "Moscow",
		Items: []model.CreateOrderItem{{
			MenuItemID:      itemID,
			Quantity:        1,
			ObservedVersion: 1,
		}},
	}

	txManager.EXPECT().Do(ctx, gomock.Any()).DoAndReturn(
		func(txCtx context.Context, fn func(context.Context) error) error {
			return fn(txCtx)
		},
	)
	repository.EXPECT().GetMenuItemsForUpdate(ctx, []uuid.UUID{itemID}).Return(
		[]model.MenuItem{{
			ID:                itemID,
			Name:              "Margherita",
			PriceMinor:        59000,
			Currency:          "RUB",
			IsActive:          true,
			IsAvailable:       true,
			AvailableQuantity: 10,
			Version:           2,
		}},
		nil,
	)

	_, err := service.CreateOrder(ctx, input)
	domainError, ok := errs.As(err)
	if !ok || domainError.Kind != errs.KindFailedPrecondition {
		t.Fatalf("expected failed precondition, got %v", err)
	}
}
