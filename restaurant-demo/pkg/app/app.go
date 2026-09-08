package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	partnerapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/api/partner/v1"
	kitchenclient "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/client/grpc/kitchen/v1"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/config"
	restaurantrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/restaurant"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/service/restaurant"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

type App struct {
	pool              *pgxpool.Pool
	grpcServer        *grpc.Server
	kitchenConnection *grpc.ClientConn
	grpcHealth        *health.Server
	grpcAddress       string
	shutdownTimeout   time.Duration
	logger            *slog.Logger
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

	kitchenConnection, err := grpc.NewClient(
		cfg.KitchenGRPCTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create Kitchen gRPC connection: %w", err)
	}

	repository := restaurantrepo.New(pool, trmpgx.DefaultCtxGetter)
	kitchenClient := kitchenclient.New(
		partnerpb.NewKitchenPartnerServiceClient(kitchenConnection),
		cfg.KitchenAPIKey,
	)
	service := restaurant.New(repository, txManager, kitchenClient, logger)
	grpcServer := grpc.NewServer()
	partnerpb.RegisterRestaurantServiceServer(
		grpcServer,
		partnerapi.New(cfg.RestaurantID, service, logger),
	)

	grpcHealth := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, grpcHealth)
	grpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	return &App{
		pool:              pool,
		grpcServer:        grpcServer,
		kitchenConnection: kitchenConnection,
		grpcHealth:        grpcHealth,
		grpcAddress:       cfg.GRPCAddress,
		shutdownTimeout:   cfg.ShutdownTimeout,
		logger:            logger,
	}, nil
}

func (a *App) Close() {
	_ = a.kitchenConnection.Close()
	a.pool.Close()
}
