package order

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

const (
	defaultPageSize = 50
	maxPageSize     = 100
)

type pageCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

func (s *Service) List(
	ctx context.Context,
	restaurantID uuid.UUID,
	status model.OrderStatus,
	pageSize int,
	pageToken string,
) (model.OrderPage, error) {
	if restaurantID == uuid.Nil {
		return model.OrderPage{}, errs.InvalidArgument("restaurant id is required")
	}
	if !validStatus(status) {
		return model.OrderPage{}, errs.InvalidArgument("order status is invalid")
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if pageSize < 1 || pageSize > maxPageSize {
		return model.OrderPage{}, errs.InvalidArgument("page size must be between 1 and 100")
	}

	filter := model.OrderFilter{
		RestaurantID: restaurantID,
		Status:       status,
		Limit:        pageSize + 1,
	}
	if pageToken != "" {
		cursor, err := decodePageToken(pageToken)
		if err != nil {
			return model.OrderPage{}, errs.InvalidArgument("page token is invalid")
		}
		filter.AfterCreatedAt = &cursor.CreatedAt
		filter.AfterID = &cursor.ID
	}

	orders, err := s.orders.ListOrders(ctx, filter)
	if err != nil {
		return model.OrderPage{}, err
	}

	page := model.OrderPage{Orders: orders}
	if len(orders) > pageSize {
		page.Orders = orders[:pageSize]
		last := page.Orders[len(page.Orders)-1]
		page.NextPageToken, err = encodePageToken(pageCursor{CreatedAt: last.CreatedAt, ID: last.ID})
		if err != nil {
			return model.OrderPage{}, errs.Internal("encode page token", err)
		}
	}

	return page, nil
}

func encodePageToken(cursor pageCursor) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodePageToken(token string) (pageCursor, error) {
	data, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return pageCursor{}, err
	}

	var cursor pageCursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return pageCursor{}, err
	}
	if cursor.CreatedAt.IsZero() || cursor.ID == uuid.Nil {
		return pageCursor{}, errs.InvalidArgument("empty page cursor")
	}

	return cursor, nil
}
