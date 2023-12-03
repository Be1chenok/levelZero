package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Stan     StanConfig
}

type ServerConfig struct {
	Host        string
	Port        int
	RequestTime time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type StanConfig struct {
	Host        string
	Port        int
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
			ServerConfig{
				Host:        viper.GetString("SERVER_HOST"),
				Port:        viper.GetInt("SERVER_PORT"),
				RequestTime: viper.GetDuration("REQUEST_TIME") * time.Second,
			},
			PostgresConfig{
				Host:     viper.GetString("PG_HOST"),
				Port:     viper.GetInt("PG_PORT"),
				Username: viper.GetString("PG_USER"),
				Password: viper.GetString("PG_PASS"),
				DBName:   viper.GetString("PG_BASE"),
				SSLMode:  viper.GetString("PG_SSL_MODE"),
			},
			StanConfig{
				Host:        viper.GetString("NATS_HOST"),
				Port:        viper.GetInt("NATS_PORT"),
				ClusterID:   viper.GetString("NATS_CLUSTER_ID"),
				ClientID:    viper.GetString("NATS_CLIENT_ID"),
				Subject:     viper.GetString("NATS_SUBJECT"),
				AckWait:     viper.GetDuration("NATS_ACK_WAIT") * time.Second,
				DurableName: viper.GetString("NATS_DURABLE_NAME"),
			},
		},
		nil
}
