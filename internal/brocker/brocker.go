package brocker

import (
	"github.com/Be1chenok/levelZero/internal/config"
	"github.com/nats-io/stan.go"
)

func New(conf *config.Config) (stan.Conn, error) {
	conn, err := stan.Connect(
		conf.Stan.ClusterID,
		conf.Stan.ClientID,
		stan.NatsURL(conf.Stan.Host+":"+conf.Stan.Port))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
