package converter

import (
	"net/url"
	"strings"

	openapi "github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/shared/pkg/openapi/kitchen/v1"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/model"
)

func CreateOrderToModel(source *openapi.CreateOrderRequest) model.CreateOrder {
	items := make([]model.CreateOrderItem, 0, len(source.Items))
	for _, item := range source.Items {
		items = append(items, model.CreateOrderItem{
			MenuItemID:      item.MenuItemID,
			Quantity:        item.Quantity,
			ObservedVersion: item.ObservedVersion,
		})
	}

	return model.CreateOrder{
		RestaurantID:    source.RestaurantID,
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		Items:           items,
	}
}

func RestaurantsToOpenAPI(source []model.Restaurant) *openapi.RestaurantList {
	restaurants := make([]openapi.Restaurant, 0, len(source))
	for _, restaurant := range source {
		restaurants = append(restaurants, openapi.Restaurant{
			ID:              restaurant.ID,
			Name:            restaurant.Name,
			AcceptingOrders: restaurant.AcceptingOrders,
		})
	}

	return &openapi.RestaurantList{Restaurants: restaurants}
}

func MenuToOpenAPI(source model.RestaurantMenu) *openapi.Menu {
	items := make([]openapi.MenuItem, 0, len(source.Items))
	for _, item := range source.Items {
		converted := openapi.MenuItem{
			ID:          item.ID,
			Name:        item.Name,
			PriceMinor:  item.PriceMinor,
			Currency:    item.Currency,
			IsAvailable: item.IsAvailable,
			Version:     item.Version,
		}
		if item.Description != "" {
			converted.Description = openapi.NewOptString(item.Description)
		}
		items = append(items, converted)
	}

	return &openapi.Menu{
		RestaurantID:    source.Restaurant.ID,
		AcceptingOrders: source.Restaurant.AcceptingOrders,
		Items:           items,
	}
}

func OrderToOpenAPI(source model.Order) *openapi.Order {
	items := make([]openapi.OrderItem, 0, len(source.Items))
	for _, item := range source.Items {
		items = append(items, openapi.OrderItem{
			MenuItemID:     item.MenuItemID,
			Name:           item.Name,
			UnitPriceMinor: item.UnitPriceMinor,
			Quantity:       item.Quantity,
			LineTotalMinor: item.LineTotalPrice,
		})
	}

	result := &openapi.Order{
		ID:              source.ID,
		RestaurantID:    source.RestaurantID,
		Status:          openapi.OrderStatus(source.Status),
		CustomerName:    source.CustomerName,
		CustomerPhone:   source.CustomerPhone,
		DeliveryAddress: source.DeliveryAddress,
		Items:           items,
		TotalMinor:      source.TotalPrice,
		Currency:        source.Currency,
		CreatedAt:       source.CreatedAt,
		UpdatedAt:       source.UpdatedAt,
	}
	if source.RejectionReason != "" {
		result.RejectionReason = openapi.NewOptString(source.RejectionReason)
	}

	return result
}

func Error(err error, statusCode int) *openapi.Error {
	code := errs.CodeInternal
	title := "Internal server error"
	detail := "an unexpected error occurred"

	if domainErr, ok := errs.As(err); ok {
		code = domainErr.Code
		title = problemTitle(domainErr.Kind)
		detail = domainErr.Message
	}

	return &openapi.Error{
		Type:   url.URL{Scheme: "https", Host: "avito-kitchen.local", Path: "/errors/" + strings.ToLower(code)},
		Title:  title,
		Status: statusCode,
		Code:   code,
		Detail: openapi.NewOptString(detail),
	}
}

func problemTitle(kind errs.Kind) string {
	switch kind {
	case errs.KindInvalidArgument:
		return "Invalid request"
	case errs.KindNotFound:
		return "Resource not found"
	case errs.KindAlreadyExists:
		return "Resource already exists"
	case errs.KindFailedPrecondition:
		return "Request cannot be completed"
	default:
		return "Internal server error"
	}
}
