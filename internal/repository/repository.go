package repository

import (
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	postgres.Order
}

func New(db *sqlx.DB) *Repository {
	return &Repository{}
}
