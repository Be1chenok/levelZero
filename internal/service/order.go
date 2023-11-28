package service

import (
	"context"

	"github.com/Be1chenok/levelZero/internal/repository/postgres"
)

type Order interface {
	FindByUID(ctx context.Context, id int)
}

type order struct {
	postgresOrder postgres.Order
}

func NewOrder(postgresOrder postgres.Order) Order {
	return &order{
		postgresOrder: postgresOrder,
	}
}

func (o order) FindByUID(ctx context.Context, id int) {

}
