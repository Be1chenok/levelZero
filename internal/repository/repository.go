package repository

import (
	"database/sql"

	"github.com/Be1chenok/levelZero/internal/config"
	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	"github.com/Be1chenok/levelZero/internal/repository/subscriber"
	"github.com/nats-io/stan.go"
)

type Repository struct {
	Subscriber    subscriber.Subscriber
	PostgresOrder postgres.Order
	CacheOrder    cache.Cache
}

func New(conf *config.Config, db *sql.DB, sc stan.Conn) *Repository {
	postgresOrder := postgres.NewOrderRepo(db)
	cacheOrder := cache.New()

	return &Repository{
		Subscriber:    subscriber.New(conf, sc, postgresOrder, cacheOrder),
		PostgresOrder: postgresOrder,
		CacheOrder:    cacheOrder,
	}
}
