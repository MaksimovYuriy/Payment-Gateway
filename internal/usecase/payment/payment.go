package payment

import (
	"context"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apperr"
	"payment_gateway/internal/lib/bankprocessor"
	"payment_gateway/internal/repo"
	"payment_gateway/internal/usecase"
	"time"

	"github.com/go-playground/validator/v10"
)

type UseCase struct {
	paymentRepo   repo.Payment
	pAttemptRepo  repo.PaymentAttempt
	bankRepo      repo.Bank
	merchantRepo  repo.Merchant
	bankProcessor bankprocessor.BankProcessor
	uow           UnitOfWork
	validate      *validator.Validate
}

func NewUseCase(
	paymentRepo repo.Payment,
	pAttemptRepo repo.PaymentAttempt,
	bankRepo repo.Bank,
	merchantRepo repo.Merchant,
	bankProcessor bankprocessor.BankProcessor,
	uow UnitOfWork,
	vd *validator.Validate,
) *UseCase {
	return &UseCase{
		paymentRepo:   paymentRepo,
		pAttemptRepo:  pAttemptRepo,
		bankRepo:      bankRepo,
		merchantRepo:  merchantRepo,
		bankProcessor: bankProcessor,
		uow:           uow,
		validate:      vd,
	}
}

var _ usecase.Payment = (*UseCase)(nil)

func (uc *UseCase) Create(ctx context.Context, p *entity.Payment, bankId int64) (*entity.Payment, error) {
	if err := uc.validateCreatePayment(p); err != nil {
		return nil, err
	}
	if err := uc.validateBankAndMerchant(ctx, bankId, p.MerchantId); err != nil {
		return nil, err
	}

	var attempt entity.PaymentAttempt
	if err := uc.createProcessingAttempt(ctx, p, bankId, &attempt); err != nil {
		return nil, err
	}

	if err := uc.processPayment(ctx, p, &attempt, bankId); err != nil {
		return nil, err
	}

	if err := uc.savePaymentResult(ctx, p, &attempt); err != nil {
		return nil, err
	}

	return p, nil
}

func (uc *UseCase) validateCreatePayment(p *entity.Payment) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if err := uc.validate.Struct(p); err != nil {
		return apperr.InvalidInput(apperr.MessageInvalidInput)
	}
	return nil
}

func (uc *UseCase) validateBankAndMerchant(ctx context.Context, bankId, merchantId int64) error {
	bank, err := uc.bankRepo.GetById(ctx, bankId)
	if err != nil {
		return err
	}
	if !bank.IsActive {
		return apperr.InvalidInput("bank is inactive")
	}
	merch, err := uc.merchantRepo.GetById(ctx, merchantId)
	if err != nil {
		return err
	}
	if !merch.IsActive {
		return apperr.InvalidInput("merchant is incative")
	}
	return nil
}

func (uc *UseCase) createProcessingAttempt(
	ctx context.Context,
	p *entity.Payment,
	bankId int64,
	attempt *entity.PaymentAttempt,
) error {
	return uc.uow.Do(ctx, func(ctx context.Context, repos Repositories) error {
		if err := repos.Payment.Create(ctx, p); err != nil {
			return err
		}

		p.SetStatusProcessing()
		if err := repos.Payment.Update(ctx, p); err != nil {
			return err
		}

		*attempt = entity.PaymentAttempt{
			PaymentId: p.Id,
			BankId:    bankId,
		}
		if err := uc.validate.Struct(attempt); err != nil {
			return apperr.InvalidInput(apperr.MessageInvalidInput)
		}

		if err := repos.PaymentAttempt.Create(ctx, attempt); err != nil {
			return err
		}

		attempt.SetStatusProcessing()
		return repos.PaymentAttempt.Update(ctx, attempt)
	})
}

func (uc *UseCase) processPayment(
	ctx context.Context,
	p *entity.Payment,
	attempt *entity.PaymentAttempt,
	bankId int64,
) error {
	processCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan bankprocessor.Result, 1)
	go func() {
		resultCh <- uc.bankProcessor.Process(processCtx, p, bankId)
	}()

	select {
	case result := <-resultCh:
		uc.applyBankResult(p, attempt, bankId, result)
	case <-time.After(5 * time.Second):
		cancel()
		attempt.Status = entity.PAttemptStatusTimeout
		p.Status = entity.PaymentStatusFailed
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (uc *UseCase) savePaymentResult(
	ctx context.Context,
	p *entity.Payment,
	attempt *entity.PaymentAttempt,
) error {
	return uc.uow.Do(ctx, func(ctx context.Context, repos Repositories) error {
		if err := repos.PaymentAttempt.Update(ctx, attempt); err != nil {
			return err
		}

		return repos.Payment.Update(ctx, p)
	})
}

func (uc *UseCase) GetById(ctx context.Context, id int64) (*entity.Payment, error) {
	return uc.paymentRepo.GetById(ctx, id)
}

func (uc *UseCase) List(ctx context.Context) ([]*entity.Payment, error) {
	return uc.paymentRepo.List(ctx)
}

func (uc *UseCase) applyBankResult(
	p *entity.Payment,
	attempt *entity.PaymentAttempt,
	bankId int64,
	result bankprocessor.Result,
) {
	if result.Success {
		p.Status = entity.PaymentStatusCompleted
		p.SuccessBankId = &bankId
		attempt.Status = entity.PAttemptStatusSucceeded
		attempt.ExternalPaymentId = stringPtr(result.ExternalPaymentId)
		return
	}

	p.Status = entity.PaymentStatusFailed
	attempt.Status = entity.PAttemptStatusFailed
	attempt.ErrorCode = stringPtr(result.ErrorCode)
	attempt.ErrorMessage = stringPtr(result.ErrorMessage)
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
