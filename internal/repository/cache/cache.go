package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/Be1chenok/levelZero/internal/domain"
)

type Cache interface {
	Get(key string) (domain.Order, bool)
	Set(key string, value domain.Order, expiration time.Duration) error
}

type cache struct {
	mutex sync.RWMutex
	data  map[string]cacheEntry
}

type cacheEntry struct {
	Value      domain.Order
	Expiration time.Time
}

func New() Cache {
	return &cache{
		data: make(map[string]cacheEntry),
	}
}

func (c *cache) Get(key string) (domain.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, ok := c.data[key]
	if !ok {
		return entry.Value, false
	}

	if time.Now().After(entry.Expiration) {
		c.mutex.Lock()
		delete(c.data, key)
		return entry.Value, false
	}

	return entry.Value, true
}

func (c *cache) Set(key string, value domain.Order, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if entry, ok := c.data[key]; ok == true {
		if !time.Now().After(entry.Expiration) {
			return fmt.Errorf("already exixts")
		}
	}

	c.data[key] = cacheEntry{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}
	return nil
}
