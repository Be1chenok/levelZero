package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   serverConfig
	Postgres postgresConfig
	Stan     stanConfig
}

type serverConfig struct {
	Host        string
	Port        string
	RequestTime time.Duration
}

type postgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type stanConfig struct {
	Host        string
	Port        string
	ClusterID   string
	ClientID    string
	Subject     string
	AckWait     time.Duration
	DurableName string
}

func Init() (*Config, error) {
	viper.SetConfigFile("../../.env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
			serverConfig{
				Host:        viper.GetString("SERVER_HOST"),
				Port:        viper.GetString("SERVER_PORT"),
				RequestTime: viper.GetDuration("REQUEST_TIME") * time.Second,
			},
			postgresConfig{
				Host:     viper.GetString("PG_HOST"),
				Port:     viper.GetString("PG_PORT"),
				Username: viper.GetString("PG_USER"),
				Password: viper.GetString("PG_PASS"),
				DBName:   viper.GetString("PG_BASE"),
				SSLMode:  viper.GetString("PG_SSL_MODE"),
			},
			stanConfig{
				Host:        viper.GetString("NATS_HOST"),
				Port:        viper.GetString("NATS_PORT"),
				ClusterID:   viper.GetString("NATS_CLUSTER_ID"),
				ClientID:    viper.GetString("NATS_CLIENT_ID"),
				Subject:     viper.GetString("NATS_SUBJECT"),
				AckWait:     viper.GetDuration("NATS_ACK_WAIT") * time.Second,
				DurableName: viper.GetString("NATS_DURABLE_NAME"),
			},
		},
		nil
}
