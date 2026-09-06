package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (c *Client) GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error) {
	response, err := c.client.GetMenu(ctx, &partnerpb.GetMenuRequest{RestaurantId: restaurantID.String()})
	if err != nil {
		return nil, clientError("get restaurant menu", err)
	}

	items, err := menuItemsFromProto(restaurantID, response.GetItems())
	if err != nil {
		return nil, errs.Internal("convert restaurant menu", err)
	}
	return items, nil
}
