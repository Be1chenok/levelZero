package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Be1chenok/levelZero/internal/config"
	_ "github.com/lib/pq"
)

func New(conf *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.Username,
		conf.Postgres.Password,
		conf.Postgres.DBName,
		conf.Postgres.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
