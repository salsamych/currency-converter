package provider

import "time"

// CurrencyRate представляет курс одной валюты к рублю
type CurrencyRate struct {
	CharCode string  // буквенный код (USD, EUR)
	Nominal  int     // номинал (обычно 1)
	Value    float64 // курс в рублях за Nominal единиц
	Name     string  // название валюты
}

// Rate содержит все курсы на определённую дату
type Rate struct {
	Date  time.Time               // дата курсов
	Rates map[string]CurrencyRate // ключ - буквенный код валюты
}

// RateProvider умеет получать текущие курсы
type RateProvider interface {
	FetchRates() (*Rate, error)
}
