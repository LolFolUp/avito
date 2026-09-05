package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"

	kitchenapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/kitchen/v1"
	partnerapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/partner/v1"
	kitchenrepo "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/repository/kitchen"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/catalog"
	orderservice "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/service/order"
)

const (
	defaultDatabaseURL   = "postgres://kitchen:kitchen@localhost:5432/kitchen?sslmode=disable"
	defaultHTTPAddress   = ":8080"
	defaultGRPCAddress   = ":9090"
	defaultPingTimeout   = 5 * time.Second
	defaultShutdown      = 10 * time.Second
	defaultHeaderTimeout = 5 * time.Second
)

type Config struct {
	DatabaseURL     string
	HTTPAddress     string
	GRPCAddress     string
	PingTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type App struct {
	pool            *pgxpool.Pool
	httpServer      *http.Server
	grpcServer      *grpc.Server
	grpcHealth      *health.Server
	grpcAddress     string
	shutdownTimeout time.Duration
	logger          *slog.Logger
}

func ConfigFromEnv() Config {
	return Config{
		DatabaseURL:     envOrDefault("DATABASE_URL", defaultDatabaseURL),
		HTTPAddress:     envOrDefault("KITCHEN_HTTP_ADDR", defaultHTTPAddress),
		GRPCAddress:     envOrDefault("KITCHEN_GRPC_ADDR", defaultGRPCAddress),
		PingTimeout:     defaultPingTimeout,
		ShutdownTimeout: defaultShutdown,
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

	repository := kitchenrepo.New(pool, trmpgx.DefaultCtxGetter)
	catalogService := catalog.New(repository, repository, txManager)
	orderService := orderservice.New(repository, repository, repository, txManager)

	restAPI := kitchenapi.New(catalogService, orderService, pool, logger)
	httpHandler, err := openapi.NewServer(restAPI)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("create HTTP handler: %w", err)
	}

	grpcServer := grpc.NewServer()
	partnerpb.RegisterPartnerServiceServer(grpcServer, partnerapi.New(catalogService, orderService, logger))
	grpcHealth := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, grpcHealth)
	grpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = defaultShutdown
	}

	return &App{
		pool: pool,
		httpServer: &http.Server{
			Addr:              cfg.HTTPAddress,
			Handler:           httpHandler,
			ReadHeaderTimeout: defaultHeaderTimeout,
		},
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

	errorsChannel := make(chan error, 2)
	go func() {
		a.logger.Info("HTTP server started", "address", a.httpServer.Addr)
		errorsChannel <- a.httpServer.ListenAndServe()
	}()
	go func() {
		a.logger.Info("gRPC server started", "address", a.grpcAddress)
		errorsChannel <- a.grpcServer.Serve(listener)
	}()

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-errorsChannel:
		if errors.Is(runErr, http.ErrServerClosed) || errors.Is(runErr, grpc.ErrServerStopped) {
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
	a.pool.Close()
}

func (a *App) shutdown() {
	a.grpcHealth.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("shutdown HTTP server", "error", err)
	}

	stopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-shutdownCtx.Done():
		a.grpcServer.Stop()
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
