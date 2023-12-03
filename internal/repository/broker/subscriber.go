package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Be1chenok/levelZero/internal/config"
	"github.com/Be1chenok/levelZero/internal/domain"
	"github.com/Be1chenok/levelZero/internal/repository/cache"
	"github.com/Be1chenok/levelZero/internal/repository/postgres"
	appLogger "github.com/Be1chenok/levelZero/logger"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type Subscriber interface {
	Subscribe(wg *sync.WaitGroup, ctx context.Context) error
	UnSubscribe() error
}

type subscriber struct {
	conf          *config.Config
	logger        appLogger.Logger
	sub           stan.Subscription
	sc            stan.Conn
	postgresOrder postgres.Order
	cache         cache.Cache
}

func NewSubscriber(conf *config.Config, logger appLogger.Logger, sc stan.Conn, postgresOrder postgres.Order, cache cache.Cache) Subscriber {
	return &subscriber{
		conf:          conf,
		sc:            sc,
		postgresOrder: postgresOrder,
		cache:         cache,
		logger:        logger.With(zap.String("component", "subscriber")),
	}
}

func (s *subscriber) Subscribe(wg *sync.WaitGroup, ctx context.Context) error {
	var err error
	s.sub, err = s.sc.Subscribe(s.conf.Stan.Subject, func(msg *stan.Msg) {
		wg.Add(1)
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			s.logger.Info("message received")
			if err := s.messageHandler(msg.Data, ctx); err != nil {
				s.logger.Errorf("failed to handle message: %v", err)
			}
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

func (s subscriber) messageHandler(data []byte, ctx context.Context) error {
	var receivedOrder domain.Order
	if err := json.Unmarshal(data, &receivedOrder); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if err := s.postgresOrder.AddOrder(ctx, receivedOrder); err != nil {
		return fmt.Errorf("failed to add order in data base: %w", err)
	}
	s.logger.Infof("order has been added to database: %s", receivedOrder.UID)

	if err := s.cache.Set(receivedOrder.UID, receivedOrder); err != nil {
		return fmt.Errorf("failed to add order in cache :%w", err)
	}
	s.logger.Infof("order has been added to cache: %s", receivedOrder.UID)

	return nil
}
