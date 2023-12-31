package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Be1chenok/levelZero/internal/config"
	appHandler "github.com/Be1chenok/levelZero/internal/delivery/http/handler"
	appServer "github.com/Be1chenok/levelZero/internal/delivery/http/server"
	appRepository "github.com/Be1chenok/levelZero/internal/repository"
	appBroker "github.com/Be1chenok/levelZero/internal/repository/broker"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	appService "github.com/Be1chenok/levelZero/internal/service"
	appLogger "github.com/Be1chenok/levelZero/logger"
	"go.uber.org/zap"
)

func Run() {
	logger, err := appLogger.NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatalf("failed to sync logger: %v", err)
		}
	}()
	appLog := logger.With(zap.String("component", "app"))

	conf, err := config.Init()
	if err != nil {
		appLog.Fatalf("failed to init config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	postgres, err := postgres.New(conf, ctx)
	if err != nil {
		appLog.Fatalf("failed to connect database: %v", err)
	}

	broker, err := appBroker.New(conf)
	if err != nil {
		appLog.Fatalf("failed to connect nats-streaming server: %v", err)
	}

	repository := appRepository.New(conf, logger, postgres, broker)
	service := appService.New(repository, logger)
	handler := appHandler.New(conf, service)
	server := appServer.New(conf, handler.InitRoutes())

	if err := repository.CacheOrder.LoadToCache(ctx); err != nil {
		log.Fatalf("failed to load cache: %v", err)
	}
	cancel()

	wg := sync.WaitGroup{}
	ctx, cancel = context.WithCancel(context.Background())

	if err := repository.Broker.Subscribe(&wg, ctx); err != nil {
		appLog.Fatalf("failed to subscribe to channel")
	}

	go func() {
		if err := server.Start(); err != nil {
			appLog.Fatalf("failed to start server: %v", err)
		}
	}()

	appLog.Infof("server is running on port %v", conf.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit
	cancel()
	wg.Wait()

	appLog.Info("shuthing down")

	if err := repository.Broker.UnSubscribe(); err != nil {
		appLog.Fatalf("failed to unsubscribe channel: %v", err)
	}

	if err := broker.Close(); err != nil {
		appLog.Fatalf("failed to close nats-streaming server connection: %v", err)
	}

	if err := server.Shuthdown(context.Background()); err != nil {
		appLog.Fatalf("failed to shut down server: %v", err)
	}

	if err := postgres.Close(); err != nil {
		appLog.Fatalf("failed to close database connection: %v", err)
	}
}
