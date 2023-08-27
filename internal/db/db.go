package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	DbHost string `envconfig:"DB_HOST"`
	DbPort string `envconfig:"DB_PORT"`
	DbName string `envconfig:"DB_NAME"`
	DbUser string `envconfig:"DB_USER"`
	DbPass string `envconfig:"DB_PASS"`
}

func CreateConnection(config DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DbHost, config.DbUser, config.DbPass, config.DbName, config.DbPort)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
