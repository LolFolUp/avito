package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	partnerapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/api/partner/v1"
	restaurantrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/restaurant"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/service/restaurant"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

const (
	defaultDatabaseURL     = "postgres://restaurant:restaurant@localhost:15433/restaurant?sslmode=disable"
	defaultGRPCAddress     = ":9090"
	defaultRestaurantID    = "00000000-0000-0000-0000-000000000001"
	defaultCurrency        = "RUB"
	defaultPingTimeout     = 5 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

type Config struct {
	DatabaseURL     string
	GRPCAddress     string
	RestaurantID    uuid.UUID
	Currency        string
	PingTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type App struct {
	pool            *pgxpool.Pool
	grpcServer      *grpc.Server
	grpcHealth      *health.Server
	grpcAddress     string
	shutdownTimeout time.Duration
	logger          *slog.Logger
}

func ConfigFromEnv() (Config, error) {
	restaurantID, err := uuid.Parse(envOrDefault("RESTAURANT_ID", defaultRestaurantID))
	if err != nil {
		return Config{}, fmt.Errorf("parse restaurant id: %w", err)
	}
	currency := strings.ToUpper(strings.TrimSpace(envOrDefault("RESTAURANT_CURRENCY", defaultCurrency)))
	if len(currency) != 3 {
		return Config{}, fmt.Errorf("restaurant currency must contain three letters")
	}

	return Config{
		DatabaseURL:     envOrDefault("DATABASE_URL", defaultDatabaseURL),
		GRPCAddress:     envOrDefault("RESTAURANT_GRPC_ADDR", defaultGRPCAddress),
		RestaurantID:    restaurantID,
		Currency:        currency,
		PingTimeout:     defaultPingTimeout,
		ShutdownTimeout: defaultShutdownTimeout,
	}, nil
}

func New(ctx context.Context, cfg Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	pingTimeout := cfg.PingTimeout
	if pingTimeout <= 0 {
		pingTimeout = defaultPingTimeout
	}
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
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

	repository := restaurantrepo.New(pool, trmpgx.DefaultCtxGetter)
	service := restaurant.New(cfg.RestaurantID, cfg.Currency, repository, txManager, logger)
	grpcServer := grpc.NewServer()
	partnerpb.RegisterRestaurantServiceServer(grpcServer, partnerapi.New(service, logger))

	grpcHealth := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, grpcHealth)
	grpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = defaultShutdownTimeout
	}
	return &App{
		pool:            pool,
		grpcServer:      grpcServer,
		grpcHealth:      grpcHealth,
		grpcAddress:     cfg.GRPCAddress,
		shutdownTimeout: shutdownTimeout,
		logger:          logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	listener, err := net.Listen("tcp", a.grpcAddress)
	if err != nil {
		return fmt.Errorf("listen gRPC: %w", err)
	}

	serveError := make(chan error, 1)
	go func() {
		a.logger.Info("restaurant gRPC server started", "address", a.grpcAddress)
		serveError <- a.grpcServer.Serve(listener)
	}()

	select {
	case <-ctx.Done():
	case err := <-serveError:
		return fmt.Errorf("serve gRPC: %w", err)
	}

	a.grpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	stopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()

	shutdownTimeout := time.NewTimer(a.shutdownTimeout)
	defer shutdownTimeout.Stop()
	select {
	case <-stopped:
	case <-shutdownTimeout.C:
		a.grpcServer.Stop()
	}
	return nil
}

func (a *App) Close() {
	a.pool.Close()
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
