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
	defaultGRPCAddress          = ":9090"
	defaultRestaurantGRPCRoutes = "00000000-0000-0000-0000-000000000001=localhost:9091"
	defaultPingTimeout          = 5 * time.Second
	defaultShutdownTimeout      = 10 * time.Second
)

type Config struct {
	DatabaseURL          string
	HTTPAddress          string
	GRPCAddress          string
	RestaurantGRPCRoutes map[uuid.UUID]string
	RestaurantAPIKeys    map[uuid.UUID]string
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
	restaurantAPIKeys, err := parseRestaurantAPIKeys(os.Getenv("RESTAURANT_API_KEYS"))
	if err != nil {
		return Config{}, err
	}
	for restaurantID := range restaurantRoutes {
		if _, exists := restaurantAPIKeys[restaurantID]; !exists {
			return Config{}, fmt.Errorf("API key for restaurant %s is required", restaurantID)
		}
	}
	for restaurantID := range restaurantAPIKeys {
		if _, exists := restaurantRoutes[restaurantID]; !exists {
			return Config{}, fmt.Errorf("gRPC route for restaurant %s is required", restaurantID)
		}
	}

	cfg := Config{
		DatabaseURL:          envOrDefault("DATABASE_URL", defaultDatabaseURL),
		HTTPAddress:          envOrDefault("KITCHEN_HTTP_ADDR", defaultHTTPAddress),
		GRPCAddress:          envOrDefault("KITCHEN_GRPC_ADDR", defaultGRPCAddress),
		RestaurantGRPCRoutes: restaurantRoutes,
		RestaurantAPIKeys:    restaurantAPIKeys,
		PingTimeout:          defaultPingTimeout,
		ShutdownTimeout:      defaultShutdownTimeout,
	}

	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if strings.TrimSpace(cfg.HTTPAddress) == "" {
		return Config{}, fmt.Errorf("KITCHEN_HTTP_ADDR is required")
	}
	if strings.TrimSpace(cfg.GRPCAddress) == "" {
		return Config{}, fmt.Errorf("KITCHEN_GRPC_ADDR is required")
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

func parseRestaurantAPIKeys(value string) (map[uuid.UUID]string, error) {
	if strings.TrimSpace(value) == "" {
		return nil, fmt.Errorf("RESTAURANT_API_KEYS is required")
	}

	keys := make(map[uuid.UUID]string)
	usedKeys := make(map[string]struct{})
	for _, rawEntry := range strings.Split(value, ",") {
		parts := strings.SplitN(strings.TrimSpace(rawEntry), "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("invalid RESTAURANT_API_KEYS entry")
		}

		restaurantID, err := uuid.Parse(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("parse RESTAURANT_API_KEYS restaurant id: %w", err)
		}
		if _, exists := keys[restaurantID]; exists {
			return nil, fmt.Errorf("duplicate API key configuration for restaurant %s", restaurantID)
		}

		apiKey := strings.TrimSpace(parts[1])
		if _, exists := usedKeys[apiKey]; exists {
			return nil, fmt.Errorf("restaurant API keys must be unique")
		}
		keys[restaurantID] = apiKey
		usedKeys[apiKey] = struct{}{}
	}

	return keys, nil
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
