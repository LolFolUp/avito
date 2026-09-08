package v1

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const apiKeyMetadataName = "x-api-key"

type authenticatedRestaurantIDKey struct{}

// AuthUnaryInterceptor проверяет API key заведения и сохраняет связанный UUID в context.
func AuthUnaryInterceptor(apiKeys map[uuid.UUID]string) grpc.UnaryServerInterceptor {
	restaurantByAPIKey := make(map[string]uuid.UUID, len(apiKeys))
	for restaurantID, apiKey := range apiKeys {
		restaurantByAPIKey[apiKey] = restaurantID
	}

	return func(
		ctx context.Context,
		request any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if strings.HasPrefix(info.FullMethod, "/grpc.health.v1.Health/") {
			return handler(ctx, request)
		}

		incomingMetadata, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "partner API key is required")
		}
		values := incomingMetadata.Get(apiKeyMetadataName)
		if len(values) != 1 {
			return nil, status.Error(codes.Unauthenticated, "partner API key is required")
		}

		restaurantID, ok := restaurantByAPIKey[values[0]]
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "partner API key is invalid")
		}

		authenticatedContext := context.WithValue(ctx, authenticatedRestaurantIDKey{}, restaurantID)
		return handler(authenticatedContext, request)
	}
}

func authenticatedRestaurantID(ctx context.Context) (uuid.UUID, bool) {
	restaurantID, ok := ctx.Value(authenticatedRestaurantIDKey{}).(uuid.UUID)
	return restaurantID, ok
}
