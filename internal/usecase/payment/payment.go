package payment

import (
	"context"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/bankprocessor"
	"payment_gateway/internal/repo"
	"payment_gateway/internal/usecase"
	"time"
)

type UseCase struct {
	paymentRepo   repo.Payment
	pAttemptRepo  repo.PaymentAttempt
	bankProcessor bankprocessor.BankProcessor
}

func NewUseCase(
	paymentRepo repo.Payment,
	pAttemptRepo repo.PaymentAttempt,
	bankProcessor bankprocessor.BankProcessor,
) *UseCase {
	return &UseCase{
		paymentRepo:   paymentRepo,
		pAttemptRepo:  pAttemptRepo,
		bankProcessor: bankProcessor,
	}
}

var _ usecase.Payment = (*UseCase)(nil)

func (uc *UseCase) Create(ctx context.Context, p *entity.Payment, bankId int64) (*entity.Payment, error) {
	if err := uc.paymentRepo.Create(ctx, p); err != nil {
		return nil, err
	}

	p.Status = entity.PaymentStatusProcessing
	if err := uc.paymentRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	attempt := entity.PaymentAttempt{
		PaymentId: p.Id,
		BankId:    bankId,
	}
	if err := uc.pAttemptRepo.Create(ctx, &attempt); err != nil {
		return nil, err
	}

	attempt.Status = entity.PAttemptStatusProcessing
	if err := uc.pAttemptRepo.Update(ctx, &attempt); err != nil {
		return nil, err
	}

	processCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan bankprocessor.Result, 1)
	go func() {
		resultCh <- uc.bankProcessor.Process(processCtx, p, bankId)
	}()

	select {
	case result := <-resultCh:
		uc.applyBankResult(p, &attempt, bankId, result)
	case <-time.After(5 * time.Second):
		cancel()
		attempt.Status = entity.PAttemptStatusTimeout
		p.Status = entity.PaymentStatusFailed
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if err := uc.pAttemptRepo.Update(ctx, &attempt); err != nil {
		return nil, err
	}
	if err := uc.paymentRepo.Update(ctx, p); err != nil {
		return nil, err
	}

	return p, nil
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
