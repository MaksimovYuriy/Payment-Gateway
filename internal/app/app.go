package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"payment_gateway/internal/config"
	"payment_gateway/internal/controller/restapi"
	"payment_gateway/internal/lib/sl"
	"payment_gateway/internal/storage/postgres"
)

func Run() error {
	appCtx, appStop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appStop()

	cfg, err := config.Load(config.ConfigPath)
	if err != nil {
		return err
	}

	logger := sl.New(cfg.App.Env)
	logger.Info("Logger started")

	database, err := postgres.New(cfg.DB)
	if err != nil {
		return err
	}
	defer database.Close()
	logger.Info("Database connection started")

	router := restapi.NewRouter()
	server := restapi.NewServer(cfg.HTTP, router)

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info(fmt.Sprintf("Server started on %s", cfg.HTTP.Port))
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-appCtx.Done():
		logger.Info("Payment-Gateway shutting down")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	return nil
}
