package subscriber

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Be1chenok/levelZero/internal/config"
	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	appLogger "github.com/Be1chenok/levelZero/logger"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Subscriber interface {
	Subscribe() error
	UnSubscribe() error
}

type subscriber struct {
	conf   *config.Config
	logger appLogger.Logger
	sub    stan.Subscription
	sc     stan.Conn
	db     postgres.Order
	cache  cache.Cache
}

func New(conf *config.Config, logger appLogger.Logger, sc stan.Conn, db postgres.Order, cache cache.Cache) Subscriber {
	return &subscriber{
		conf:   conf,
		sc:     sc,
		db:     db,
		cache:  cache,
		logger: logger.With(zap.String("component", "subscriber")),
	}
}

func (s *subscriber) Subscribe() error {
	var err error

	s.sub, err = s.sc.Subscribe(s.conf.Stan.Subject, func(msg *stan.Msg) {
		s.logger.Info("message received")
		if err := s.messageHandler(msg.Data); err == nil {
			if err := msg.Ack(); err != nil {
				s.logger.Infof("failed to acknowledge message: %v\n", err)
			}
		}
	},
		stan.AckWait(s.conf.Stan.AckWait),
		stan.DurableName(s.conf.Stan.DurableName),
		stan.SetManualAckMode(),
		stan.MaxInflight(15))
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	s.logger.Infof("subscribe succesful")

	return nil
}

func (s *subscriber) UnSubscribe() error {
	if s.sub != nil {
		if err := s.sub.Unsubscribe(); err != nil {
			return fmt.Errorf("failed to unsubscribe: %w", err)
		}
	}

	return nil
}

func (s subscriber) messageHandler(data []byte) error {
	var receivedOrder domain.Order
	if err := json.Unmarshal(data, &receivedOrder); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if err := s.db.AddOrder(context.Background(), receivedOrder); err != nil {
		return fmt.Errorf("failed to add order in data base: %w", err)
	}
	s.logger.Infof("order has been added to database: %s", receivedOrder.OrderUID)

	if err := s.cache.Set(receivedOrder.OrderUID, receivedOrder); err != nil {
		return fmt.Errorf("failed to add order in cache :%w", err)
	}
	s.logger.Infof("order has been added to cache: %s", receivedOrder.OrderUID)

	return nil
}
