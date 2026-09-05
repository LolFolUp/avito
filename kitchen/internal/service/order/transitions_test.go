package order

import (
	"testing"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name string
		from model.OrderStatus
		to   model.OrderStatus
		want bool
	}{
		{"accept pending order", model.OrderStatusPendingConfirmation, model.OrderStatusAccepted, true},
		{"reject pending order", model.OrderStatusPendingConfirmation, model.OrderStatusRejected, true},
		{"start preparation", model.OrderStatusAccepted, model.OrderStatusPreparing, true},
		{"skip preparation", model.OrderStatusAccepted, model.OrderStatusReady, false},
		{"change terminal order", model.OrderStatusDelivered, model.OrderStatusPreparing, false},
		{"accept rejected order", model.OrderStatusRejected, model.OrderStatusAccepted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canTransition(tt.from, tt.to); got != tt.want {
				t.Fatalf("canTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
