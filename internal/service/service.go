package service

import (
	"github.com/Be1chenok/levelZero/internal/repository"
)

type Service struct {
	Order
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Order: NewOrder(repo.PostgresOrder, repo.CacheOrder, repo.Subscriber),
	}
}
