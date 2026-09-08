package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/api/converter"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func (a *API) UpdateOrderStatus(
	ctx context.Context,
	request *partnerpb.UpdateOrderStatusRequest,
) (*partnerpb.UpdateOrderStatusResponse, error) {
	authenticatedID, ok := authenticatedRestaurantID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "partner authentication is required")
	}

	requestedID, orderID, expectedStatus, newStatus, err := converter.UpdateOrderStatusFromProto(request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if requestedID != authenticatedID {
		return nil, status.Error(codes.PermissionDenied, "restaurant id does not match API key")
	}

	order, err := a.orders.UpdateStatus(ctx, authenticatedID, orderID, expectedStatus, newStatus, "")
	if err != nil {
		a.logError("update order status", err)
		return nil, grpcError(err)
	}

	return &partnerpb.UpdateOrderStatusResponse{
		OrderId:   order.ID.String(),
		Status:    converter.OrderStatusToProto(order.Status),
		UpdatedAt: timestamppb.New(order.UpdatedAt),
	}, nil
}
