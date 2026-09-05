package order

import "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"

var allowedTransitions = map[model.OrderStatus]map[model.OrderStatus]struct{}{
	model.OrderStatusPendingConfirmation: {
		model.OrderStatusAccepted: {},
		model.OrderStatusRejected: {},
	},
	model.OrderStatusAccepted: {
		model.OrderStatusPreparing: {},
	},
	model.OrderStatusPreparing: {
		model.OrderStatusReady: {},
	},
	model.OrderStatusReady: {
		model.OrderStatusDelivering: {},
	},
	model.OrderStatusDelivering: {
		model.OrderStatusDelivered: {},
	},
}

func validStatus(status model.OrderStatus) bool {
	switch status {
	case model.OrderStatusPendingConfirmation,
		model.OrderStatusAccepted,
		model.OrderStatusPreparing,
		model.OrderStatusReady,
		model.OrderStatusDelivering,
		model.OrderStatusDelivered,
		model.OrderStatusRejected:
		return true
	default:
		return false
	}
}

func canTransition(from, to model.OrderStatus) bool {
	_, allowed := allowedTransitions[from][to]
	return allowed
}
