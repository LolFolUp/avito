package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (c *Client) UpdateOrderStatus(
	ctx context.Context,
	restaurantID uuid.UUID,
	orderID uuid.UUID,
	expectedStatus model.OrderStatus,
	newStatus model.OrderStatus,
	reason string,
) (model.Order, error) {
	response, err := c.client.UpdateOrderStatus(ctx, &partnerpb.UpdateOrderStatusRequest{
		RestaurantId:   restaurantID.String(),
		OrderId:        orderID.String(),
		ExpectedStatus: statusToProto(expectedStatus),
		NewStatus:      statusToProto(newStatus),
		Reason:         reason,
	})
	if err != nil {
		return model.Order{}, clientError(err)
	}
	if response.GetOrder() == nil {
		return model.Order{}, &errs.Error{Kind: errs.KindInternal, Message: "partner returned an empty order"}
	}

	order, err := orderFromProto(response.GetOrder())
	if err != nil {
		return model.Order{}, &errs.Error{Kind: errs.KindInternal, Message: "decode order", Cause: err}
	}
	return order, nil
}
