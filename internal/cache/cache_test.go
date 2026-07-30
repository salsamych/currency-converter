package cache_test

import (
	"testing"
	"time"

	"currency-converter/internal/cache"
	"currency-converter/internal/provider"
)

func TestCacheSetAndGet(t *testing.T) {
	c := cache.NewCache(100 * time.Millisecond)

	rates := &provider.Rate{
		Date: time.Now(),
		Rates: map[string]provider.CurrencyRate{
			"USD": {CharCode: "USD", Nominal: 1, Value: 88.07},
		},
	}

	// Кеш должен быть пустым
	if got := c.Get(); got != nil {
		t.Error("cache should be empty initially")
	}

	// Сохраняем данные
	c.Set(rates)

	// Данные должны быть доступны
	if got := c.Get(); got == nil {
		t.Error("cache should return data after Set")
	}

	// Ждём истечения TTL
	time.Sleep(150 * time.Millisecond)

	// Кеш должен быть снова пустым
	if got := c.Get(); got != nil {
		t.Error("cache should be empty after TTL")
	}
}

func TestCacheMultipleSets(t *testing.T) {
	c := cache.NewCache(1 * time.Second)

	rates1 := &provider.Rate{
		Date: time.Now(),
		Rates: map[string]provider.CurrencyRate{
			"USD": {CharCode: "USD", Value: 88.07},
		},
	}

	rates2 := &provider.Rate{
		Date: time.Now(),
		Rates: map[string]provider.CurrencyRate{
			"USD": {CharCode: "USD", Value: 90.00},
		},
	}

	c.Set(rates1)
	c.Set(rates2)

	got := c.Get()
	if got.Rates["USD"].Value != 90.00 {
		t.Errorf("expected updated value 90.00, got %v", got.Rates["USD"].Value)
	}
}

func TestCacheConcurrency(t *testing.T) {
	c := cache.NewCache(1 * time.Second)

	rates := &provider.Rate{
		Date: time.Now(),
		Rates: map[string]provider.CurrencyRate{
			"USD": {CharCode: "USD", Value: 88.07},
		},
	}

	// Тест на конкурентный доступ
	done := make(chan bool)

	go func() {
		for i := 0; i < 100; i++ {
			c.Set(rates)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			c.Get()
		}
		done <- true
	}()

	<-done
	<-done

	// Проверяем, что данные не повреждены
	if got := c.Get(); got == nil {
		t.Error("cache should have data after concurrent access")
	}
}
