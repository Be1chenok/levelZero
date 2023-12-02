package service

import (
	"context"
	"fmt"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	"github.com/Be1chenok/levelZero/internal/repository/subscriber"
	appLogger "github.com/Be1chenok/levelZero/logger"
	"go.uber.org/zap"
)

type Order interface {
	FindByUID(ctx context.Context, orderUID string) (domain.Order, error)
	UnSubscribeToChannel() error
	SubscribeToChannel() error
	LoadToCache() error
}

type order struct {
	postgresOrder postgres.Order
	cacheOrder    cache.Cache
	subscriber    subscriber.Subscriber
	logger        appLogger.Logger
}

func NewOrder(postgresOrder postgres.Order, cacheOrder cache.Cache, subscriber subscriber.Subscriber, logger appLogger.Logger) Order {
	return &order{
		postgresOrder: postgresOrder,
		cacheOrder:    cacheOrder,
		subscriber:    subscriber,
		logger:        logger.With(zap.String("component", "service-order")),
	}
}

func (o order) SubscribeToChannel() error {
	if err := o.subscriber.Subscribe(); err != nil {
		return err
	}

	return nil
}

func (o order) UnSubscribeToChannel() error {
	if err := o.subscriber.UnSubscribe(); err != nil {
		return err
	}

	return nil
}

func (o order) LoadToCache() error {
	orders, err := o.postgresOrder.FindAllOrders()
	if err != nil {
		return fmt.Errorf("failed to find all orders: %w", err)
	}
	if len(orders) != 0 {
		for i := range orders {
			if err := o.cacheOrder.Set(orders[i].OrderUID, orders[i]); err != nil {
				return fmt.Errorf("filed to set data: %w", err)
			}
		}
	}
	o.logger.Infof("load to cache %v orders", len(orders))

	return nil
}

func (o order) FindByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	order, ok := o.cacheOrder.Get(orderUID)
	if ok {
		return order, nil
	}

	order, err := o.postgresOrder.FindOrderByUID(ctx, orderUID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to find order by UID: %w", err)
	}

	delivery, err := o.postgresOrder.FindDeliveryByOrderUID(ctx, orderUID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to find delivery by orderUID: %w", err)
	}

	payment, err := o.postgresOrder.FindPaymentByOrderUID(ctx, orderUID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to find payment by orderUID: %w", err)
	}

	items, err := o.postgresOrder.FindItemsByOrderUID(ctx, orderUID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to find items by orderUID: %w", err)
	}

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	return order, nil
}
