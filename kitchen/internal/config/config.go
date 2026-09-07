package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	defaultDatabaseURL          = "postgres://kitchen:kitchen@localhost:15432/kitchen?sslmode=disable"
	defaultHTTPAddress          = ":8080"
	defaultRestaurantGRPCRoutes = "00000000-0000-0000-0000-000000000001=localhost:9090"
	defaultPingTimeout          = 5 * time.Second
	defaultShutdownTimeout      = 10 * time.Second
)

type Config struct {
	DatabaseURL          string
	HTTPAddress          string
	RestaurantGRPCRoutes map[uuid.UUID]string
	PingTimeout          time.Duration
	ShutdownTimeout      time.Duration
}

func Load() (Config, error) {
	restaurantRoutes, err := parseRestaurantRoutes(
		envOrDefault("RESTAURANT_GRPC_ROUTES", defaultRestaurantGRPCRoutes),
	)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		DatabaseURL:          envOrDefault("DATABASE_URL", defaultDatabaseURL),
		HTTPAddress:          envOrDefault("KITCHEN_HTTP_ADDR", defaultHTTPAddress),
		RestaurantGRPCRoutes: restaurantRoutes,
		PingTimeout:          defaultPingTimeout,
		ShutdownTimeout:      defaultShutdownTimeout,
	}

	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if strings.TrimSpace(cfg.HTTPAddress) == "" {
		return Config{}, fmt.Errorf("KITCHEN_HTTP_ADDR is required")
	}
	return cfg, nil
}

func parseRestaurantRoutes(value string) (map[uuid.UUID]string, error) {
	routes := make(map[uuid.UUID]string)
	for _, rawRoute := range strings.Split(value, ",") {
		parts := strings.SplitN(strings.TrimSpace(rawRoute), "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("invalid RESTAURANT_GRPC_ROUTES entry %q", rawRoute)
		}

		restaurantID, err := uuid.Parse(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("parse restaurant route id %q: %w", parts[0], err)
		}
		if _, exists := routes[restaurantID]; exists {
			return nil, fmt.Errorf("duplicate restaurant route %s", restaurantID)
		}
		routes[restaurantID] = strings.TrimSpace(parts[1])
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("RESTAURANT_GRPC_ROUTES is required")
	}
	return routes, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
