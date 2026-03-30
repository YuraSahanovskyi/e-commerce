package db

import (
	"e-commerce/internal/config"
	"e-commerce/internal/logger"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_HOST"),
		config.GetEnv("DB_PORT"),
		config.GetEnv("DB_NAME"),
	)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		logger.Log.Error("migration init error:" + err.Error())
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		logger.Log.Error("migration error:" + err.Error())
	}

	logger.Log.Info("migrations applied")
}
