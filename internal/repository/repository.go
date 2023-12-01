package repository

import (
	"database/sql"

	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
)

type Repository struct {
	PostgresOrder postgres.Order
	CacheOrder    cache.Cache
}

func New(db *sql.DB) *Repository {
	postgresOrder := postgres.NewOrderRepo(db)
	return &Repository{
		PostgresOrder: postgresOrder,
		CacheOrder:    cache.New(),
	}
}
