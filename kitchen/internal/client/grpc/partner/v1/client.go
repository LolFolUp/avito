package v1

import (
	"google.golang.org/grpc"

	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

type Client struct {
	client partnerpb.RestaurantServiceClient
}

func New(connection grpc.ClientConnInterface) *Client {
	return &Client{client: partnerpb.NewRestaurantServiceClient(connection)}
}
