package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kitchenapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/kitchen/v1"
	partnerclient "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/client/grpc/partner/v1"
	kitchenrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/kitchen"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/catalog"
	orderservice "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/order"
	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
)

const (
	defaultDatabaseURL           = "postgres://kitchen:kitchen@localhost:15432/kitchen?sslmode=disable"
	defaultHTTPAddress           = ":8080"
	defaultRestaurantGRPCAddress = "localhost:9090"
	defaultPingTimeout           = 5 * time.Second
	defaultShutdown              = 10 * time.Second
	defaultHeaderTimeout         = 5 * time.Second
)

type Config struct {
	DatabaseURL           string
	HTTPAddress           string
	RestaurantGRPCAddress string
	PingTimeout           time.Duration
	ShutdownTimeout       time.Duration
}

type App struct {
	pool              *pgxpool.Pool
	httpServer        *http.Server
	partnerConnection *grpc.ClientConn
	shutdownTimeout   time.Duration
	logger            *slog.Logger
}

func ConfigFromEnv() Config {
	return Config{
		DatabaseURL:           envOrDefault("DATABASE_URL", defaultDatabaseURL),
		HTTPAddress:           envOrDefault("KITCHEN_HTTP_ADDR", defaultHTTPAddress),
		RestaurantGRPCAddress: envOrDefault("RESTAURANT_GRPC_ADDR", defaultRestaurantGRPCAddress),
		PingTimeout:           defaultPingTimeout,
		ShutdownTimeout:       defaultShutdown,
	}
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

	partnerConnection, err := grpc.NewClient(
		cfg.RestaurantGRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create restaurant gRPC connection: %w", err)
	}
	partnerClient := partnerclient.New(partnerConnection)
	repository := kitchenrepo.New(pool, trmpgx.DefaultCtxGetter)
	catalogService := catalog.New(repository, partnerClient)
	orderService := orderservice.New(repository, partnerClient, repository, txManager)

	restAPI := kitchenapi.New(catalogService, orderService, pool, logger)
	httpHandler, err := openapi.NewServer(restAPI)
	if err != nil {
		_ = partnerConnection.Close()
		pool.Close()
		return nil, fmt.Errorf("create HTTP handler: %w", err)
	}

	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = defaultShutdown
	}

	return &App{
		pool:              pool,
		partnerConnection: partnerConnection,
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddress,
			Handler:           httpHandler,
			ReadHeaderTimeout: defaultHeaderTimeout,
		},
		shutdownTimeout: shutdownTimeout,
		logger:          logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errorsChannel := make(chan error, 1)
	go func() {
		a.logger.Info("HTTP server started", "address", a.httpServer.Addr)
		errorsChannel <- a.httpServer.ListenAndServe()
	}()

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-errorsChannel:
		if errors.Is(runErr, http.ErrServerClosed) {
			runErr = nil
		}
	}

	a.shutdown()
	if runErr != nil {
		return fmt.Errorf("serve: %w", runErr)
	}
	return nil
}

func (a *App) Close() {
	_ = a.partnerConnection.Close()
	a.pool.Close()
}

func (a *App) shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("shutdown HTTP server", "error", err)
	}

}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
