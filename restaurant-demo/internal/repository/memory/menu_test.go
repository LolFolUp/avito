package memory

import (
	"context"
	"testing"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/model"
)

func TestItemsAvailable(t *testing.T) {
	repository := New([]model.MenuItem{
		{ExternalID: "available", IsActive: true, IsAvailable: true},
		{ExternalID: "disabled", IsActive: false, IsAvailable: true},
	})

	tests := []struct {
		name string
		ids  []string
		want bool
	}{
		{"all available", []string{"available"}, true},
		{"inactive", []string{"disabled"}, false},
		{"missing", []string{"missing"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.ItemsAvailable(context.Background(), tt.ids)
			if err != nil {
				t.Fatalf("items available: %v", err)
			}
			if got != tt.want {
				t.Fatalf("available = %v, want %v", got, tt.want)
			}
		})
	}
}
