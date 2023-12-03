package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	appLogger "github.com/Be1chenok/levelZero/logger"
	"go.uber.org/zap"
)

var ErrAlreadyExists = errors.New("already exixts")

type Cache interface {
	LoadToCache(ctx context.Context) error
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}) error
}

type cache struct {
	logger        appLogger.Logger
	postgresOrder postgres.Order
	mutex         sync.RWMutex
	data          map[string]interface{}
}

func New(postgresOrder postgres.Order, logger appLogger.Logger) Cache {
	return &cache{
		data:          make(map[string]interface{}),
		postgresOrder: postgresOrder,
		logger:        logger.With(zap.String("component", "cache")),
	}
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}

	return value, true
}

func (c *cache) Set(key string, value interface{}) error {
	c.mutex.RLock()
	if _, ok := c.data[key]; ok == true {
		return ErrAlreadyExists
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	c.data[key] = value
	c.mutex.Unlock()

	return nil
}

func (c *cache) LoadToCache(ctx context.Context) error {
	orders, err := c.postgresOrder.FindAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to find all orders: %w", err)
	}

	c.logger.Info("loading cache")

	for _, order := range orders {
		if err := c.Set(order.UID, order); err != nil {
			return fmt.Errorf("filed to set data: %w", err)
		}
	}
	c.logger.Infof("loaded to cache %v orders", len(orders))

	return nil
}
