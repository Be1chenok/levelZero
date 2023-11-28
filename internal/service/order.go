package service

import (
	"context"
	"fmt"

	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
)

type Order interface {
	FindByUID(ctx context.Context, orderUID string) (domain.Order, error)
}

type order struct {
	postgresOrder postgres.Order
}

func NewOrder(postgresOrder postgres.Order) Order {
	return &order{
		postgresOrder: postgresOrder,
	}
}

func (o order) FindByUID(ctx context.Context, orderUID string) (domain.Order, error) {
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
