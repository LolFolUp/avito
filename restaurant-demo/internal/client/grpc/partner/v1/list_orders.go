package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (c *Client) ListOrders(
	ctx context.Context,
	restaurantID uuid.UUID,
	status model.OrderStatus,
	pageSize int,
	pageToken string,
) (model.OrderPage, error) {
	response, err := c.client.ListOrders(ctx, &partnerpb.ListOrdersRequest{
		RestaurantId: restaurantID.String(),
		Status:       statusToProto(status),
		PageSize:     int32(pageSize),
		PageToken:    pageToken,
	})
	if err != nil {
		return model.OrderPage{}, clientError(err)
	}

	orders := make([]model.Order, 0, len(response.GetOrders()))
	for _, source := range response.GetOrders() {
		order, err := orderFromProto(source)
		if err != nil {
			return model.OrderPage{}, &errs.Error{Kind: errs.KindInternal, Message: "decode order", Cause: err}
		}
		orders = append(orders, order)
	}

	return model.OrderPage{Orders: orders, NextPageToken: response.GetNextPageToken()}, nil
}
