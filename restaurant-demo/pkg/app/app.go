package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	partnerclient "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/client/grpc/partner/v1"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/repository/memory"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/service/restaurant"
)

const (
	defaultKitchenAddress = "localhost:9090"
	defaultRestaurantID   = "00000000-0000-0000-0000-000000000001"
	defaultPollInterval   = 2 * time.Second
)

type Config struct {
	KitchenAddress string
	RestaurantID   uuid.UUID
	PollInterval   time.Duration
}

type App struct {
	connection   *grpc.ClientConn
	service      *restaurant.Service
	pollInterval time.Duration
}

func ConfigFromEnv() (Config, error) {
	restaurantID, err := uuid.Parse(envOrDefault("RESTAURANT_ID", defaultRestaurantID))
	if err != nil {
		return Config{}, fmt.Errorf("parse restaurant id: %w", err)
	}

	pollInterval, err := time.ParseDuration(envOrDefault("RESTAURANT_POLL_INTERVAL", defaultPollInterval.String()))
	if err != nil {
		return Config{}, fmt.Errorf("parse restaurant poll interval: %w", err)
	}
	if pollInterval <= 0 {
		return Config{}, fmt.Errorf("restaurant poll interval must be positive")
	}

	return Config{
		KitchenAddress: envOrDefault("KITCHEN_GRPC_ADDR", defaultKitchenAddress),
		RestaurantID:   restaurantID,
		PollInterval:   pollInterval,
	}, nil
}

func New(cfg Config) (*App, error) {
	connection, err := grpc.NewClient(
		cfg.KitchenAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("create Kitchen gRPC connection: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	menuRepository := memory.New(demoMenu())
	client := partnerclient.New(connection)
	service := restaurant.New(cfg.RestaurantID, client, menuRepository, logger)

	return &App{connection: connection, service: service, pollInterval: cfg.PollInterval}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	return a.service.Run(ctx, a.pollInterval)
}

func (a *App) Close() error {
	return a.connection.Close()
}

func demoMenu() []model.MenuItem {
	return []model.MenuItem{
		{
			ExternalID:  "pizza-margherita",
			Name:        "Margherita",
			Description: "Tomato sauce and mozzarella",
			PriceMinor:  59000,
			Currency:    "RUB",
			IsActive:    true,
			IsAvailable: true,
		},
		{
			ExternalID:  "pasta-carbonara",
			Name:        "Carbonara",
			Description: "Pasta, bacon and cream sauce",
			PriceMinor:  65000,
			Currency:    "RUB",
			IsActive:    true,
			IsAvailable: true,
		},
		{
			ExternalID:  "tiramisu",
			Name:        "Tiramisu",
			Description: "Classic Italian dessert",
			PriceMinor:  35000,
			Currency:    "RUB",
			IsActive:    true,
			IsAvailable: false,
		},
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
