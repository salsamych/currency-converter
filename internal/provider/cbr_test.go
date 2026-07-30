package provider_test

import (
	"testing"

	"currency-converter/internal/provider"
)

func TestCBRProviderIntegration(t *testing.T) {
	p := provider.NewCBRProvider()

	rates, err := p.FetchRates()
	if err != nil {
		t.Fatalf("failed to fetch rates from CBR: %v", err)
	}

	// Проверяем, что курсы не пустые
	if len(rates.Rates) == 0 {
		t.Error("received empty rates")
	}

	// Проверяем наличие основных валют
	requiredCurrencies := []string{"USD", "EUR", "RUB"}
	for _, code := range requiredCurrencies {
		if _, exists := rates.Rates[code]; !exists {
			t.Errorf("required currency %s not found", code)
		}
	}

	// Проверяем, что дата не нулевая
	if rates.Date.IsZero() {
		t.Error("date should not be zero")
	}

	// Проверяем курс рубля
	if rubRate, ok := rates.Rates["RUB"]; ok {
		if rubRate.Value != 1.0 {
			t.Errorf("RUB rate should be 1.0, got %f", rubRate.Value)
		}
	}

	// Проверяем, что все курсы положительные
	for code, rate := range rates.Rates {
		if rate.Value <= 0 {
			t.Errorf("rate for %s should be positive, got %f", code, rate.Value)
		}
		if rate.Nominal <= 0 {
			t.Errorf("nominal for %s should be positive, got %d", code, rate.Nominal)
		}
	}
}
