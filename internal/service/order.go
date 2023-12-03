package service

import (
	"context"
	"fmt"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	appLogger "github.com/Be1chenok/levelZero/logger"
	"go.uber.org/zap"
)

type Order interface {
	FindByUID(ctx context.Context, orderUID string) (domain.Order, error)
}

type order struct {
	postgresOrder postgres.Order
	cacheOrder    cache.Cache
	logger        appLogger.Logger
}

func NewOrder(postgresOrder postgres.Order, cacheOrder cache.Cache, logger appLogger.Logger) Order {
	return &order{
		postgresOrder: postgresOrder,
		cacheOrder:    cacheOrder,
		logger:        logger.With(zap.String("component", "service-order")),
	}
}

func (o order) FindByUID(ctx context.Context, orderUID string) (domain.Order, error) {
	cachedOrder, ok := o.cacheOrder.Get(orderUID)
	if ok {
		arr, ok := cachedOrder.(domain.Order)
		if !ok {
			o.logger.Error("failed to convert data")
		}
		return arr, nil
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

	if err := o.cacheOrder.Set(order.UID, order); err != nil {
		o.logger.Errorf("failed to add order %s to cache: %v", order.UID, err)
	}

	o.logger.Infof("order %s added to cache", order.UID)

	return order, nil
}
