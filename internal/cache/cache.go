package cache

import (
	"sync"
	"time"

	"currency-converter/internal/provider"
)

// Cache хранит курсы валют в памяти с TTL
type Cache struct {
	mu        sync.RWMutex
	rates     *provider.Rate
	ttl       time.Duration
	updatedAt time.Time
}

// NewCache создаёт новый кеш с указанным временем жизни
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl: ttl,
	}
}

// Get возвращает курсы из кеша, если они не устарели
func (c *Cache) Get() *provider.Rate {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.rates != nil && time.Since(c.updatedAt) < c.ttl {
		return c.rates
	}
	return nil
}

// Set сохраняет курсы в кеш
func (c *Cache) Set(rates *provider.Rate) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rates = rates
	c.updatedAt = time.Now()
}
