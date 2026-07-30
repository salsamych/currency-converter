package converter_test

import (
	"testing"
	"time"

	"currency-converter/internal/cache"
	"currency-converter/internal/converter"
	"currency-converter/internal/provider"
)

// mockProvider имитирует ЦБ РФ
type mockProvider struct {
	rates *provider.Rate
	err   error
}

func (m *mockProvider) FetchRates() (*provider.Rate, error) {
	return m.rates, m.err
}

func TestConvert(t *testing.T) {
	// Создаём тестовые курсы как в реальном ЦБ
	rates := &provider.Rate{
		Date: time.Now(),
		Rates: map[string]provider.CurrencyRate{
			"USD": {
				CharCode: "USD",
				Nominal:  1,
				Value:    75.00, // 1 USD = 75.00 RUB
				Name:     "Доллар США",
			},
			"EUR": {
				CharCode: "EUR",
				Nominal:  1,
				Value:    90.00, // 1 EUR = 90.00 RUB
				Name:     "Евро",
			},
			"RUB": {
				CharCode: "RUB",
				Nominal:  1,
				Value:    1.0,
				Name:     "Российский рубль",
			},
			"JPY": {
				CharCode: "JPY",
				Nominal:  100,   // курс за 100 иен!
				Value:    50.00, // 100 JPY = 50.00 RUB
				Name:     "Японская иена",
			},
		},
	}

	mock := &mockProvider{rates: rates, err: nil}
	cache := cache.NewCache(10 * time.Minute)
	conv := converter.NewConverter(mock, cache)

	tests := []struct {
		name    string
		from    string
		to      string
		amount  float64
		want    float64
		wantErr bool
	}{
		{
			name:    "USD to RUB",
			from:    "USD",
			to:      "RUB",
			amount:  100,
			want:    7500.0, // 100 * 75.00
			wantErr: false,
		},
		{
			name:    "RUB to USD",
			from:    "RUB",
			to:      "USD",
			amount:  7500,
			want:    100.0, // 7500 / 75.00
			wantErr: false,
		},
		{
			name:    "USD to EUR",
			from:    "USD",
			to:      "EUR",
			amount:  100,
			want:    83.33, // 100 * 75.00 / 90.00 ≈ 83.33
			wantErr: false,
		},
		{
			name:    "JPY to RUB (номинал 100)",
			from:    "JPY",
			to:      "RUB",
			amount:  1000,
			want:    500.0, // 1000 * (50.00 / 100) = 500.0
			wantErr: false,
		},
		{
			name:    "RUB to JPY (номинал 100)",
			from:    "RUB",
			to:      "JPY",
			amount:  500.0,
			want:    1000.0, // 500.00 / (50.00 / 100) = 1000
			wantErr: false,
		},
		{
			name:    "несуществующая валюта",
			from:    "XXX",
			to:      "RUB",
			amount:  100,
			want:    0,
			wantErr: true,
		},
		{
			name:    "отрицательная сумма",
			from:    "USD",
			to:      "RUB",
			amount:  -100,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := conv.Convert(tt.from, tt.to, tt.amount)

			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Сравниваем с точностью до 2 знаков после запятой
				diff := got - tt.want
				if diff > 0.01 || diff < -0.01 {
					t.Errorf("Convert() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetAvailableCurrencies(t *testing.T) {
	rates := &provider.Rate{
		Date: time.Now(),
		Rates: map[string]provider.CurrencyRate{
			"USD": {CharCode: "USD", Name: "Доллар США"},
			"EUR": {CharCode: "EUR", Name: "Евро"},
		},
	}

	mock := &mockProvider{rates: rates, err: nil}
	cache := cache.NewCache(10 * time.Minute)
	conv := converter.NewConverter(mock, cache)

	currencies, err := conv.GetAvailableCurrencies()
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if len(currencies) != 2 {
		t.Errorf("ожидалось 2 валюты, получено %d", len(currencies))
	}

	if currencies["USD"] != "Доллар США" {
		t.Errorf("неправильное название валюты USD: %s", currencies["USD"])
	}
}
