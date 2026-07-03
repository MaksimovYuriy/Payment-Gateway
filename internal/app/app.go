package app

import (
	"fmt"
	"log"
	"payment_gateway/internal/config"
)

func Run() error {
	cfg, err := config.Load(config.ConfigPath)
	if err != nil {
		log.Fatal("Config error: %w", err)
	}
	fmt.Println(cfg.DB.Host)

	// Логгирование

	// БД

	// Зависимости

	// HTTP сервер

	// Graceful shutdown

	return nil
}
