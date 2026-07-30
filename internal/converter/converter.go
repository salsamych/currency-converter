// Package converter реализует бизнес-логику конвертации валют.
package converter

import (
	"errors"
	"fmt"

	"currency-converter/internal/cache"
	"currency-converter/internal/provider"
)

// ErrInvalidCurrency возвращается при передаче неподдерживаемой валюты.
var ErrInvalidCurrency = errors.New("неподдерживаемая валюта")

// ErrFetchRates возвращается при ошибке получения курсов.
var ErrFetchRates = errors.New("ошибка получения курсов")

// ErrInvalidAmount возвращается при некорректной сумме.
var ErrInvalidAmount = errors.New("некорректная сумма")

// Converter выполняет конвертацию валют через рубль.
type Converter struct {
	provider provider.RateProvider
	cache    *cache.Cache
}

// NewConverter создаёт новый конвертер.
func NewConverter(p provider.RateProvider, c *cache.Cache) *Converter {
	return &Converter{
		provider: p,
		cache:    c,
	}
}

// Convert конвертирует amount из валюты from в валюту to.
func (c *Converter) Convert(from, to string, amount float64) (float64, error) {
	if from == "" || to == "" {
		return 0, ErrInvalidCurrency
	}
	if amount <= 0 {
		return 0, ErrInvalidAmount
	}

	rates, err := c.getRates()
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrFetchRates, err)
	}

	fromRate, exists := rates.Rates[from]
	if !exists {
		return 0, fmt.Errorf("%w: %s", ErrInvalidCurrency, from)
	}

	toRate, exists := rates.Rates[to]
	if !exists {
		return 0, fmt.Errorf("%w: %s", ErrInvalidCurrency, to)
	}

	rubAmount := amount * (fromRate.Value / float64(fromRate.Nominal))
	result := rubAmount / (toRate.Value / float64(toRate.Nominal))

	return result, nil
}

func (c *Converter) getRates() (*provider.Rate, error) {
	if cached := c.cache.Get(); cached != nil {
		return cached, nil
	}

	rates, err := c.provider.FetchRates()
	if err != nil {
		return nil, err
	}

	c.cache.Set(rates)
	return rates, nil
}

// GetAvailableCurrencies возвращает список доступных валют.
func (c *Converter) GetAvailableCurrencies() (map[string]string, error) {
	rates, err := c.getRates()
	if err != nil {
		return nil, err
	}

	currencies := make(map[string]string, len(rates.Rates))
	for code, rate := range rates.Rates {
		currencies[code] = rate.Name
	}

	return currencies, nil
}
