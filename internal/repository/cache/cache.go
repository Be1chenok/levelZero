package cache

import (
	"fmt"
	"sync"

	"github.com/Be1chenok/levelZero/internal/domain"
)

type Cache interface {
	Get(key string) (domain.Order, bool)
	Set(key string, value domain.Order) error
}

type cache struct {
	mutex sync.RWMutex
	data  map[string]domain.Order
}

func New() Cache {
	return &cache{
		data: make(map[string]domain.Order),
	}
}

func (c *cache) Get(key string) (domain.Order, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return value, false
	}

	return value, true
}

func (c *cache) Set(key string, value domain.Order) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.data[key]; ok == true {
		return fmt.Errorf("already exixts")
	}

	c.data[key] = value
	return nil
}
