// Package provider реализует получение курсов валют из внешних источников.
package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const cbrURL = "https://www.cbr-xml-daily.ru/daily_json.js"

// CBRProvider получает курсы с сайта ЦБ РФ в JSON формате.
type CBRProvider struct {
	client *http.Client
}

// NewCBRProvider создаёт новый провайдер курсов ЦБ.
func NewCBRProvider() *CBRProvider {
	return &CBRProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchRates получает актуальные курсы валют.
func (p *CBRProvider) FetchRates() (*Rate, error) {
	var response struct {
		Date   time.Time `json:"Date"`
		Valute map[string]struct {
			CharCode string  `json:"CharCode"`
			Nominal  int     `json:"Nominal"`
			Name     string  `json:"Name"`
			Value    float64 `json:"Value"`
		} `json:"Valute"`
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cbrURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к ЦБ: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус от ЦБ: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	if len(response.Valute) == 0 {
		return nil, fmt.Errorf("получен пустой список валют")
	}

	rates := make(map[string]CurrencyRate, len(response.Valute))
	for code, valute := range response.Valute {
		rates[code] = CurrencyRate{
			CharCode: valute.CharCode,
			Nominal:  valute.Nominal,
			Value:    valute.Value,
			Name:     valute.Name,
		}
	}

	rates["RUB"] = CurrencyRate{
		CharCode: "RUB",
		Nominal:  1,
		Value:    1,
		Name:     "Российский рубль",
	}

	return &Rate{
		Date:  response.Date,
		Rates: rates,
	}, nil
}
