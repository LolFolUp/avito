package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
	partnerpb "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/proto/partner/v1"
)

func MenuItemsFromProto(source []*partnerpb.MenuItem) []model.MenuItem {
	items := make([]model.MenuItem, 0, len(source))
	for _, item := range source {
		if item == nil {
			items = append(items, model.MenuItem{})
			continue
		}
		items = append(items, model.MenuItem{
			ExternalID:  item.GetExternalId(),
			Name:        item.GetName(),
			Description: item.GetDescription(),
			PriceMinor:  item.GetPriceMinor(),
			Currency:    item.GetCurrency(),
			IsActive:    item.GetIsActive(),
			IsAvailable: item.GetIsAvailable(),
		})
	}

	return items
}

func OrderStatusFromProto(source partnerpb.OrderStatus) model.OrderStatus {
	switch source {
	case partnerpb.OrderStatus_ORDER_STATUS_PENDING_CONFIRMATION:
		return model.OrderStatusPendingConfirmation
	case partnerpb.OrderStatus_ORDER_STATUS_ACCEPTED:
		return model.OrderStatusAccepted
	case partnerpb.OrderStatus_ORDER_STATUS_PREPARING:
		return model.OrderStatusPreparing
	case partnerpb.OrderStatus_ORDER_STATUS_READY:
		return model.OrderStatusReady
	case partnerpb.OrderStatus_ORDER_STATUS_DELIVERING:
		return model.OrderStatusDelivering
	case partnerpb.OrderStatus_ORDER_STATUS_DELIVERED:
		return model.OrderStatusDelivered
	case partnerpb.OrderStatus_ORDER_STATUS_REJECTED:
		return model.OrderStatusRejected
	default:
		return ""
	}
}

func OrderToProto(source model.Order) *partnerpb.Order {
	items := make([]*partnerpb.OrderItem, 0, len(source.Items))
	for _, item := range source.Items {
		items = append(items, &partnerpb.OrderItem{
			MenuItemId:     item.MenuItemID.String(),
			ExternalItemId: item.ExternalItemID,
			Name:           item.Name,
			UnitPriceMinor: item.UnitPriceMinor,
			Quantity:       int32(item.Quantity),
			LineTotalMinor: item.LineTotalMinor,
		})
	}

	return &partnerpb.Order{
		Id:              source.ID.String(),
		RestaurantId:    source.RestaurantID.String(),
		Status:          OrderStatusToProto(source.Status),
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		Items:           items,
		TotalMinor:      source.TotalMinor,
		Currency:        source.Currency,
		RejectionReason: source.RejectionReason,
		CreatedAt:       timestamp(source.CreatedAt),
		UpdatedAt:       timestamp(source.UpdatedAt),
	}
}

func OrdersToProto(source []model.Order) []*partnerpb.Order {
	result := make([]*partnerpb.Order, 0, len(source))
	for _, order := range source {
		result = append(result, OrderToProto(order))
	}
	return result
}

func OrderStatusToProto(source model.OrderStatus) partnerpb.OrderStatus {
	switch source {
	case model.OrderStatusPendingConfirmation:
		return partnerpb.OrderStatus_ORDER_STATUS_PENDING_CONFIRMATION
	case model.OrderStatusAccepted:
		return partnerpb.OrderStatus_ORDER_STATUS_ACCEPTED
	case model.OrderStatusPreparing:
		return partnerpb.OrderStatus_ORDER_STATUS_PREPARING
	case model.OrderStatusReady:
		return partnerpb.OrderStatus_ORDER_STATUS_READY
	case model.OrderStatusDelivering:
		return partnerpb.OrderStatus_ORDER_STATUS_DELIVERING
	case model.OrderStatusDelivered:
		return partnerpb.OrderStatus_ORDER_STATUS_DELIVERED
	case model.OrderStatusRejected:
		return partnerpb.OrderStatus_ORDER_STATUS_REJECTED
	default:
		return partnerpb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func timestamp(value time.Time) *timestamppb.Timestamp {
	if value.IsZero() {
		return nil
	}
	return timestamppb.New(value)
}
