package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (c *Client) SyncMenu(
	ctx context.Context,
	restaurantID uuid.UUID,
	menu []model.MenuItem,
) (model.MenuSyncResult, error) {
	response, err := c.client.SyncMenu(ctx, &partnerpb.SyncMenuRequest{
		RestaurantId: restaurantID.String(),
		Items:        menuToProto(menu),
	})
	if err != nil {
		return model.MenuSyncResult{}, clientError(err)
	}

	return model.MenuSyncResult{
		CreatedCount:     int(response.GetCreatedCount()),
		UpdatedCount:     int(response.GetUpdatedCount()),
		DeactivatedCount: int(response.GetDeactivatedCount()),
	}, nil
}
