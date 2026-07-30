package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"currency-converter/internal/api"
	"currency-converter/internal/cache"
	"currency-converter/internal/converter"
	"currency-converter/internal/provider"
)

func TestIntegrationConvertEndpoint(t *testing.T) {
	realProvider := provider.NewCBRProvider()
	cache := cache.NewCache(10 * time.Minute)
	conv := converter.NewConverter(realProvider, cache)
	handler := api.NewHandler(conv)

	server := httptest.NewServer(http.HandlerFunc(handler.Convert))
	defer server.Close()

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{"USD to RUB", "/convert?from=USD&to=RUB&amount=100", http.StatusOK},
		{"EUR to USD", "/convert?from=EUR&to=USD&amount=50", http.StatusOK},
		{"invalid currency", "/convert?from=XXX&to=RUB&amount=100", http.StatusBadRequest},
		{"negative amount", "/convert?from=USD&to=RUB&amount=-100", http.StatusBadRequest},
		{"missing params", "/convert", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.url)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {
				var result map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Fatalf("decode failed: %v", err)
				}
				if result["result"] == nil {
					t.Error("missing 'result' field")
				}
				if val, ok := result["result"].(float64); ok && val <= 0 {
					t.Error("result should be positive")
				}
			}
		})
	}
}

func TestIntegrationCurrenciesEndpoint(t *testing.T) {
	realProvider := provider.NewCBRProvider()
	cache := cache.NewCache(10 * time.Minute)
	conv := converter.NewConverter(realProvider, cache)
	handler := api.NewHandler(conv)

	server := httptest.NewServer(http.HandlerFunc(handler.Currencies))
	defer server.Close()

	resp, err := http.Get(server.URL + "/currencies")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var currencies map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		t.Fatalf("decode failed: %v", err)
	}

	if len(currencies) < 3 {
		t.Errorf("expected at least 3 currencies, got %d", len(currencies))
	}
}
