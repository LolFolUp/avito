package catalog

import (
	"testing"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func TestPrepareMenuNormalizesValues(t *testing.T) {
	items, err := prepareMenu([]model.MenuItem{{
		ExternalID:  "  pizza-1 ",
		Name:        " Margherita ",
		Description: " Classic ",
		PriceMinor:  59000,
		Currency:    " rub ",
		IsActive:    true,
		IsAvailable: true,
	}})
	if err != nil {
		t.Fatalf("prepare menu: %v", err)
	}

	got := items[0]
	if got.ExternalID != "pizza-1" || got.Name != "Margherita" || got.Description != "Classic" || got.Currency != "RUB" {
		t.Fatalf("menu item was not normalized: %#v", got)
	}
}

func TestPrepareMenuRejectsDuplicateExternalID(t *testing.T) {
	_, err := prepareMenu([]model.MenuItem{
		{ExternalID: "pizza-1", Name: "One", Currency: "RUB"},
		{ExternalID: "pizza-1", Name: "Two", Currency: "RUB"},
	})
	if err == nil {
		t.Fatal("expected duplicate external id error")
	}

	domainErr, ok := errs.As(err)
	if !ok || domainErr.Kind != errs.KindInvalidArgument {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}
