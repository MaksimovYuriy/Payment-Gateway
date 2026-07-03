package app

import (
	"log"
	"payment_gateway/internal/config"
	"payment_gateway/internal/lib/sl"
)

func Run() error {
	cfg, err := config.Load(config.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	logger := sl.New(cfg.APP.Env)
	logger.Info("Info")

	// БД

	// Зависимости

	// HTTP сервер

	// Graceful shutdown

	return nil
}
