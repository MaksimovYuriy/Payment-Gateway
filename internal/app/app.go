package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"payment_gateway/internal/config"
	"payment_gateway/internal/controller/restapi"
	"payment_gateway/internal/lib/bankprocessor"
	"payment_gateway/internal/lib/sl"
	bankrepo "payment_gateway/internal/repo/bank"
	merchantrepo "payment_gateway/internal/repo/merchant"
	paymentrepo "payment_gateway/internal/repo/payment"
	pattemptrepo "payment_gateway/internal/repo/payment_attempt"
	"payment_gateway/internal/storage/postgres"
	bankusecase "payment_gateway/internal/usecase/bank"
	merchusecase "payment_gateway/internal/usecase/merchant"
	paymentusecase "payment_gateway/internal/usecase/payment"

	"github.com/go-playground/validator/v10"
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

	validate := validator.New()

	bankProcessor := bankprocessor.NewMockProcessor(3 * time.Second)

	paymentUOW := postgres.NewUnitOfWork(database)

	bankRepo := bankrepo.NewRepo(database)
	merchRepo := merchantrepo.NewRepo(database)
	paymentRepo := paymentrepo.NewRepo(database)
	pAttemptRepo := pattemptrepo.NewRepo(database)

	bankUseCase := bankusecase.NewUseCase(bankRepo)
	merchUseCase := merchusecase.NewUseCase(merchRepo)
	paymentUseCase := paymentusecase.NewUseCase(paymentRepo, pAttemptRepo, bankProcessor, paymentUOW)

	bankController := restapi.NewBankController(bankUseCase, validate)
	merchController := restapi.NewMerchantController(merchUseCase, validate)
	paymentController := restapi.NewPaymentController(paymentUseCase, validate)

	router := restapi.NewRouter(bankController, merchController, paymentController)
	server := restapi.NewServer(cfg.HTTP, router)

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Server started", slog.String("addr", server.Addr))
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

	logger.Info("Payment-Gateway stopped")

	return nil
}
