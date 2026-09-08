package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func (c *Client) UpdateOrderStatus(ctx context.Context, input model.UpdateOrderStatus) error {
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", c.apiKey)
	response, err := c.client.UpdateOrderStatus(ctx, updateOrderStatusToProto(input))
	if err != nil {
		return clientError("update Kitchen order status", err)
	}
	if response.GetOrderId() != input.OrderID.String() {
		return errs.Internal("validate Kitchen order response", fmt.Errorf("unexpected order id"))
	}
	if response.GetStatus() != orderStatusToProto(input.NewStatus) {
		return errs.Internal("validate Kitchen order response", fmt.Errorf("unexpected order status"))
	}
	return nil
}
