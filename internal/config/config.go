package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   serverConfig
	Postgres postgresConfig
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
		},
		nil
}
