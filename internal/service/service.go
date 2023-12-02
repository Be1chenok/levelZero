package service

import (
	"github.com/Be1chenok/levelZero/internal/repository"
	appLogger "github.com/Be1chenok/levelZero/logger"
)

type Service struct {
	Order
}

func New(repo *repository.Repository, logger appLogger.Logger) *Service {
	return &Service{
		Order: NewOrder(repo.PostgresOrder, repo.CacheOrder, repo.Subscriber, logger),
	}
}
