package catalog

import (
	"context"
	"strings"
	"unicode"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func (s *Service) SyncMenu(
	ctx context.Context,
	restaurantID uuid.UUID,
	menu []model.MenuItem,
) (model.MenuSyncResult, error) {
	if restaurantID == uuid.Nil {
		return model.MenuSyncResult{}, errs.InvalidArgument("restaurant id is required")
	}

	prepared, err := prepareMenu(menu)
	if err != nil {
		return model.MenuSyncResult{}, err
	}

	var result model.MenuSyncResult
	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		restaurant, err := s.restaurants.GetRestaurantForUpdate(txCtx, restaurantID)
		if err != nil {
			return err
		}

		existing, err := s.menu.GetMenuForUpdate(txCtx, restaurantID)
		if err != nil {
			return err
		}

		existingByExternalID := make(map[string]model.MenuItem, len(existing))
		for _, item := range existing {
			existingByExternalID[item.ExternalID] = item
		}

		externalIDs := make([]string, 0, len(prepared))
		for _, incoming := range prepared {
			if incoming.Currency != restaurant.Currency {
				return errs.InvalidArgument("menu item currency must match restaurant currency")
			}

			externalIDs = append(externalIDs, incoming.ExternalID)
			current, exists := existingByExternalID[incoming.ExternalID]
			if !exists {
				incoming.ID = uuid.New()
				incoming.RestaurantID = restaurantID
				incoming.Version = 1
				if err := s.menu.UpsertMenuItem(txCtx, incoming); err != nil {
					return err
				}
				result.CreatedCount++
				continue
			}

			if menuItemEqual(current, incoming) {
				continue
			}

			incoming.ID = current.ID
			incoming.RestaurantID = restaurantID
			incoming.Version = current.Version + 1
			if err := s.menu.UpsertMenuItem(txCtx, incoming); err != nil {
				return err
			}
			result.UpdatedCount++
		}

		deactivated, err := s.menu.DeactivateMenuItemsExcept(txCtx, restaurantID, externalIDs)
		if err != nil {
			return err
		}
		result.DeactivatedCount = int(deactivated)

		return nil
	})
	if err != nil {
		return model.MenuSyncResult{}, err
	}

	return result, nil
}

func prepareMenu(menu []model.MenuItem) ([]model.MenuItem, error) {
	seen := make(map[string]struct{}, len(menu))
	prepared := make([]model.MenuItem, 0, len(menu))

	for _, item := range menu {
		item.ExternalID = strings.TrimSpace(item.ExternalID)
		item.Name = strings.TrimSpace(item.Name)
		item.Description = strings.TrimSpace(item.Description)
		item.Currency = strings.ToUpper(strings.TrimSpace(item.Currency))

		switch {
		case item.ExternalID == "" || len(item.ExternalID) > 200:
			return nil, errs.InvalidArgument("menu item external id is invalid")
		case item.Name == "" || len(item.Name) > 200:
			return nil, errs.InvalidArgument("menu item name is invalid")
		case len(item.Description) > 2000:
			return nil, errs.InvalidArgument("menu item description is too long")
		case item.PriceMinor < 0:
			return nil, errs.InvalidArgument("menu item price cannot be negative")
		case !validCurrency(item.Currency):
			return nil, errs.InvalidArgument("menu item currency is invalid")
		}

		if _, exists := seen[item.ExternalID]; exists {
			return nil, errs.InvalidArgument("menu contains duplicate external ids")
		}
		seen[item.ExternalID] = struct{}{}
		prepared = append(prepared, item)
	}

	return prepared, nil
}

func menuItemEqual(current, incoming model.MenuItem) bool {
	return current.Name == incoming.Name &&
		current.Description == incoming.Description &&
		current.PriceMinor == incoming.PriceMinor &&
		current.Currency == incoming.Currency &&
		current.IsActive == incoming.IsActive &&
		current.IsAvailable == incoming.IsAvailable
}

func validCurrency(currency string) bool {
	if len(currency) != 3 {
		return false
	}
	for _, symbol := range currency {
		if !unicode.IsUpper(symbol) || !unicode.IsLetter(symbol) {
			return false
		}
	}

	return true
}
