package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
)

type Order interface {
	FindByUID(ctx context.Context, orderUID string) (domain.Order, error)
	LoadToCache() error
}

type order struct {
	postgresOrder postgres.Order
	cacheOrder    cache.Cache
}

func NewOrder(postgresOrder postgres.Order, cacheOrder cache.Cache) Order {
	return &order{
		postgresOrder: postgresOrder,
		cacheOrder:    cacheOrder,
	}
}

func (o order) LoadToCache() error {
	orders, err := o.postgresOrder.FindAllOrders()
	if err != nil {
		return fmt.Errorf("failed to find all orders: %w", err)
	}
	if len(orders) != 0 {
		for i := range orders {
			if err := o.cacheOrder.Set(orders[i].OrderUID, orders[i], 10*time.Minute); err != nil {
				return fmt.Errorf("filed to set data: %w", err)
			}
		}
	}

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
