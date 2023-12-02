package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Be1chenok/levelZero/internal/brocker"
	"github.com/Be1chenok/levelZero/internal/config"
	appHandler "github.com/Be1chenok/levelZero/internal/delivery/http/handler"
	appServer "github.com/Be1chenok/levelZero/internal/delivery/http/server"
	appRepository "github.com/Be1chenok/levelZero/internal/repository"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	appService "github.com/Be1chenok/levelZero/internal/service"
)

func Run() {
	conf, err := config.Init()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	postgres, err := postgres.New(conf)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	brocker, err := brocker.New(conf)
	if err != nil {
		log.Fatalf("failed to connect nats-streaming server: %v", err)
	}

	repository := appRepository.New(conf, postgres, brocker)
	service := appService.New(repository)
	handler := appHandler.New(conf, service)
	server := appServer.New(conf, handler.InitRoutes())

	if err := service.LoadToCache(); err != nil {
		log.Fatalf("failed to load cache: %v", err)
	}

	if err := service.SubscribeToChannel(); err != nil {
		log.Fatalf("subscriber: %v", err)
	}

	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	log.Printf("server is running on port %v", conf.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	<-quit

	log.Print("shuthing down")

	if err := service.UnSubscribeToChannel(); err != nil {
		log.Fatalf("subscriber: %v", err)
	}

	if err := brocker.Close(); err != nil {
		log.Fatalf("failed to close nats-streaming server connection: %v", err)
	}

	if err := server.Shuthdown(context.Background()); err != nil {
		log.Fatalf("failed to shut down server: %v", err)
	}

	if err := postgres.Close(); err != nil {
		log.Fatalf("failed to close database connection: %v", err)
	}
}
