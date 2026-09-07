package partner

import (
	"context"
	"errors"
	"fmt"

	v1 "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/client/grpc/partner/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

// Resolver выбирает gRPC-клиент заведения по его UUID в Kitchen.
type Resolver struct {
	clients     map[uuid.UUID]*v1.Client
	connections []*grpc.ClientConn
}

func NewResolver(routes map[uuid.UUID]string) (*Resolver, error) {
	resolver := &Resolver{
		clients:     make(map[uuid.UUID]*v1.Client, len(routes)),
		connections: make([]*grpc.ClientConn, 0, len(routes)),
	}

	for restaurantID, target := range routes {
		connection, err := grpc.NewClient(
			target,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			_ = resolver.Close()
			return nil, fmt.Errorf("create restaurant %s gRPC connection: %w", restaurantID, err)
		}

		resolver.connections = append(resolver.connections, connection)
		resolver.clients[restaurantID] = v1.New(connection)
	}

	return resolver, nil
}

func (r *Resolver) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error) {
	client, err := r.clientFor(restaurantID)
	if err != nil {
		return nil, err
	}
	return client.GetMenu(ctx, restaurantID)
}

func (r *Resolver) CreateOrder(
	ctx context.Context,
	orderID uuid.UUID,
	input model.CreateOrder,
) (model.Order, error) {
	client, err := r.clientFor(input.RestaurantID)
	if err != nil {
		return model.Order{}, err
	}
	return client.CreateOrder(ctx, orderID, input)
}

func (r *Resolver) clientFor(restaurantID uuid.UUID) (*v1.Client, error) {
	client, exists := r.clients[restaurantID]
	if !exists {
		return nil, errs.NotFound("restaurant integration")
	}
	return client, nil
}

func (r *Resolver) Close() error {
	closeErrors := make([]error, 0, len(r.connections))
	for _, connection := range r.connections {
		if err := connection.Close(); err != nil {
			closeErrors = append(closeErrors, err)
		}
	}
	return errors.Join(closeErrors...)
}
