package repository

import (
	"database/sql"

	"github.com/Be1chenok/levelZero/internal/config"
	"github.com/Be1chenok/levelZero/internal/repository/broker"
	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	appLogger "github.com/Be1chenok/levelZero/logger"
	"github.com/nats-io/stan.go"
)

type Repository struct {
	Broker        broker.Subscriber
	PostgresOrder postgres.Order
	CacheOrder    cache.Cache
}

func New(conf *config.Config, logger appLogger.Logger, db *sql.DB, sc stan.Conn) *Repository {
	postgresOrder := postgres.NewOrderRepo(db)
	cacheOrder := cache.New(postgresOrder, logger)

	return &Repository{
		Broker:        broker.NewSubscriber(conf, logger, sc, postgresOrder, cacheOrder),
		PostgresOrder: postgresOrder,
		CacheOrder:    cacheOrder,
	}
}
