package v1

import (
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
	"google.golang.org/grpc"
)

type Client struct {
	client partnerpb.PartnerServiceClient
}

func New(connection grpc.ClientConnInterface) *Client {
	return &Client{client: partnerpb.NewPartnerServiceClient(connection)}
}
