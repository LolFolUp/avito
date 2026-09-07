//go:build integration

package restaurant_test

import (
	"context"
	"os"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	restaurantrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/restaurant"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/service/restaurant"
)

func TestCreateOrderAtomicallyReservesItem(t *testing.T) {
	databaseURL := os.Getenv("RESTAURANT_TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("RESTAURANT_TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		t.Fatalf("create transaction manager: %v", err)
	}
	repository := restaurantrepo.New(pool, trmpgx.DefaultCtxGetter)

	service := restaurant.New(repository, txManager, nil, nil)
	itemID := uuid.New()
	orderID := uuid.New()
	_, err = pool.Exec(ctx, `
		INSERT INTO menu_items (
			id, name, description, price_minor, currency,
			is_active, is_available, available_quantity, version
		) VALUES ($1, 'Item', '', 100, 'RUB', TRUE, TRUE, 1, 1)`,
		itemID,
	)
	if err != nil {
		t.Fatalf("prepare menu item: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM orders WHERE id = $1`, orderID)
		_, _ = pool.Exec(context.Background(), `DELETE FROM menu_items WHERE id = $1`, itemID)
	})

	created, err := service.CreateOrder(ctx, model.CreateOrder{
		ID:              orderID,
		CustomerName:    "Alice",
		CustomerPhone:   "+79990000000",
		DeliveryAddress: "Moscow",
		Items: []model.CreateOrderItem{{
			MenuItemID:      itemID,
			Quantity:        1,
			ObservedVersion: 1,
		}},
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	if created.Status != model.OrderStatusAccepted || created.TotalPrice != 100 {
		t.Fatalf("unexpected order: %#v", created)
	}

	var quantity int
	var available bool
	var ordersCount int
	err = pool.QueryRow(ctx, `
		SELECT available_quantity, is_available,
		       (SELECT COUNT(*) FROM orders WHERE id = $2)
		FROM menu_items
		WHERE id = $1`, itemID, orderID).Scan(&quantity, &available, &ordersCount)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	if quantity != 0 || available || ordersCount != 1 {
		t.Fatalf("reservation result: quantity=%d available=%v orders=%d", quantity, available, ordersCount)
	}
}
