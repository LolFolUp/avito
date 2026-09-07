package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	defaultDatabaseURL     = "postgres://restaurant:restaurant@localhost:15433/restaurant?sslmode=disable"
	defaultGRPCAddress     = ":9090"
	defaultRestaurantID    = "00000000-0000-0000-0000-000000000001"
	defaultPingTimeout     = 5 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

type Config struct {
	DatabaseURL     string
	GRPCAddress     string
	RestaurantID    uuid.UUID
	PingTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func Load() (Config, error) {
	restaurantID, err := uuid.Parse(envOrDefault("RESTAURANT_ID", defaultRestaurantID))
	if err != nil {
		return Config{}, fmt.Errorf("parse RESTAURANT_ID: %w", err)
	}

	cfg := Config{
		DatabaseURL:     envOrDefault("DATABASE_URL", defaultDatabaseURL),
		GRPCAddress:     envOrDefault("RESTAURANT_GRPC_ADDR", defaultGRPCAddress),
		RestaurantID:    restaurantID,
		PingTimeout:     defaultPingTimeout,
		ShutdownTimeout: defaultShutdownTimeout,
	}

	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if strings.TrimSpace(cfg.GRPCAddress) == "" {
		return Config{}, fmt.Errorf("RESTAURANT_GRPC_ADDR is required")
	}

	return cfg, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
