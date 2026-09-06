package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func (c *Client) CreateOrder(ctx context.Context, orderID uuid.UUID, input model.CreateOrder) (model.Order, error) {
	response, err := c.client.CreateRestaurantOrder(ctx, createOrderRequestToProto(orderID, input))
	if err != nil {
		return model.Order{}, clientError("create restaurant order", err)
	}

	order, err := orderFromProto(response.GetOrder())
	if err != nil {
		return model.Order{}, errs.Internal("convert restaurant order", err)
	}
	return order, nil
}
