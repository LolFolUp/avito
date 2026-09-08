package v1

import partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"

type Client struct {
	client partnerpb.KitchenPartnerServiceClient
	apiKey string
}

func New(client partnerpb.KitchenPartnerServiceClient, apiKey string) *Client {
	return &Client{client: client, apiKey: apiKey}
}
