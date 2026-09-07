//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type menuResponse struct {
	RestaurantID string     `json:"restaurant_id"`
	Items        []menuItem `json:"items"`
}

type menuItem struct {
	ID          string `json:"id"`
	IsAvailable bool   `json:"is_available"`
	Version     int64  `json:"version"`
}

type orderResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	TotalMinor int64  `json:"total_minor"`
}

func TestOrderIsAcceptedByRestaurantDemo(t *testing.T) {
	baseURL := envOrDefault("E2E_BASE_URL", "http://127.0.0.1:8080")
	restaurantID := envOrDefault("E2E_RESTAURANT_ID", "00000000-0000-0000-0000-000000000001")
	client := &http.Client{Timeout: 3 * time.Second}

	var menu menuResponse
	getJSON(t, client, baseURL+"/api/v1/restaurants/"+restaurantID+"/menu", &menu)
	var selected menuItem
	for _, item := range menu.Items {
		if item.IsAvailable {
			selected = item
			break
		}
	}
	if selected.ID == "" {
		t.Fatal("demo menu has no available items")
	}

	body := map[string]any{
		"restaurant_id":    restaurantID,
		"customer_name":    "E2E Customer",
		"customer_phone":   "+79990000000",
		"delivery_address": "Moscow",
		"items": []map[string]any{{
			"menu_item_id":     selected.ID,
			"quantity":         1,
			"observed_version": selected.Version,
		}},
	}
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal order: %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, baseURL+"/api/v1/orders", bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		message, _ := io.ReadAll(response.Body)
		t.Fatalf("create order status = %d: %s", response.StatusCode, message)
	}
	var order orderResponse
	if err := json.NewDecoder(response.Body).Decode(&order); err != nil {
		t.Fatalf("decode order: %v", err)
	}
	if order.Status != "ACCEPTED" || order.TotalMinor <= 0 {
		t.Fatalf("unexpected created order: %#v", order)
	}
	getJSON(t, client, fmt.Sprintf("%s/api/v1/orders/%s", baseURL, order.ID), &order)
	if order.Status != "ACCEPTED" {
		t.Fatalf("stored order has status %q", order.Status)
	}

	body["items"] = []map[string]any{{
		"menu_item_id":     selected.ID,
		"quantity":         1,
		"observed_version": selected.Version + 1,
	}}
	stalePayload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal stale order: %v", err)
	}
	staleRequest, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/api/v1/orders",
		bytes.NewReader(stalePayload),
	)
	if err != nil {
		t.Fatalf("create stale request: %v", err)
	}
	staleRequest.Header.Set("Content-Type", "application/json")
	staleResponse, err := client.Do(staleRequest)
	if err != nil {
		t.Fatalf("create stale order: %v", err)
	}
	defer staleResponse.Body.Close()
	if staleResponse.StatusCode != http.StatusConflict {
		message, _ := io.ReadAll(staleResponse.Body)
		t.Fatalf("stale order status = %d: %s", staleResponse.StatusCode, message)
	}
}

func getJSON(t *testing.T, client *http.Client, url string, target any) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("create GET request: %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(response.Body)
		t.Fatalf("GET %s status = %d: %s", url, response.StatusCode, message)
	}
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("decode GET %s: %v", url, err)
	}
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
