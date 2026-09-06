//go:build integration

package kitchen

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func TestRepositoryUsesTransactionFromContext(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping database: %v", err)
	}

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		t.Fatalf("create transaction manager: %v", err)
	}
	repository := New(pool, trmpgx.DefaultCtxGetter)

	restaurantID := uuid.New()
	menuItemID := uuid.New()
	orderID := uuid.New()
	_, err = pool.Exec(ctx, `INSERT INTO restaurants (id, name, currency) VALUES ($1, $2, $3)`, restaurantID, "Integration", "RUB")
	if err != nil {
		t.Fatalf("insert restaurant: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM restaurants WHERE id = $1`, restaurantID)
	})

	rollback := errors.New("force rollback")
	err = txManager.Do(ctx, func(txCtx context.Context) error {
		if err := repository.CreateOrder(txCtx, model.Order{
			ID: orderID, RestaurantID: restaurantID, Status: model.OrderStatusAccepted,
			CustomerName: "Alice", CustomerPhone: "+79990000000", DeliveryAddress: "Moscow",
			TotalMinor: 100, Currency: "RUB",
		}); err != nil {
			return err
		}
		if err := repository.CreateOrderItems(txCtx, []model.OrderItem{{
			OrderID: orderID, MenuItemID: menuItemID,
			Name: "Item", UnitPriceMinor: 100, Quantity: 1,
		}}); err != nil {
			return err
		}
		if err := repository.CreateOrderStatusHistory(txCtx, model.OrderStatusHistory{
			OrderID: orderID, Status: model.OrderStatusAccepted,
		}); err != nil {
			return err
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatalf("transaction error = %v, want rollback sentinel", err)
	}

	_, err = repository.GetOrder(ctx, orderID)
	domainErr, ok := errs.As(err)
	if !ok || domainErr.Kind != errs.KindNotFound {
		t.Fatalf("order survived rollback: %v", err)
	}
}
