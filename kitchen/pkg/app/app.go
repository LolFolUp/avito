package app

import (
	"context"
	"fmt"
	partnerclient "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/client/grpc/partner"
	"log/slog"
	"net/http"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"

	kitchenapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/kitchen/v1"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/config"
	kitchenrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/kitchen"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/catalog"
	orderservice "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/order"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

type App struct {
	pool            *pgxpool.Pool
	httpServer      *http.Server
	partnerResolver *partnerclient.Resolver
	shutdownTimeout time.Duration
	logger          *slog.Logger
}

func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	if logger == nil {
		logger = slog.Default()
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.PingTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create transaction manager: %w", err)
	}

	partnerResolver, err := partnerclient.NewResolver(cfg.RestaurantGRPCRoutes)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create restaurant gRPC resolver: %w", err)
	}

	repository := kitchenrepo.New(pool, trmpgx.DefaultCtxGetter)
	catalogService := catalog.New(repository, partnerResolver)
	orderService := orderservice.New(repository, partnerResolver, repository, txManager)

	restAPI := kitchenapi.New(catalogService, orderService, pool, logger)
	httpHandler, err := openapi.NewServer(restAPI)
	if err != nil {
		_ = partnerResolver.Close()
		pool.Close()
		return nil, fmt.Errorf("create HTTP handler: %w", err)
	}

	return &App{
		pool:            pool,
		partnerResolver: partnerResolver,
		httpServer:      newHTTPServer(cfg.HTTPAddress, httpHandler),
		shutdownTimeout: cfg.ShutdownTimeout,
		logger:          logger,
	}, nil
}

func (a *App) Close() {
	_ = a.partnerResolver.Close()
	a.pool.Close()
}
