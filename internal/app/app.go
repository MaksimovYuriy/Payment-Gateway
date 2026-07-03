package app

import (
	"log"

	"payment_gateway/internal/config"
	"payment_gateway/internal/lib/sl"
	"payment_gateway/internal/storage/postgres"
)

func Run() error {
	cfg, err := config.Load(config.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	logger := sl.New(cfg.App.Env)
	logger.Info("Info")

	// БД
	database, err := postgres.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Зависимости

	// HTTP сервер

	// Graceful shutdown

	return nil
}
