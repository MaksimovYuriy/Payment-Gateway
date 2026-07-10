package bank

import (
	"context"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apperr"
	"payment_gateway/internal/repo"
	"payment_gateway/internal/usecase"

	"github.com/go-playground/validator/v10"
)

type UseCase struct {
	bankRepo repo.Bank
	validate *validator.Validate
}

var _ usecase.Bank = (*UseCase)(nil)

func NewUseCase(br repo.Bank, vd *validator.Validate) *UseCase {
	return &UseCase{
		bankRepo: br,
		validate: vd,
	}
}

func (uc *UseCase) Registration(ctx context.Context, br *entity.Bank) (*entity.Bank, error) {
	if err := uc.validate.Struct(br); err != nil {
		return nil, apperr.InvalidInput(apperr.MessageInvalidInput)
	}
	if err := uc.bankRepo.Create(ctx, br); err != nil {
		return nil, err
	}
	return br, nil
}

func (uc *UseCase) List(ctx context.Context) ([]*entity.Bank, error) {
	banks, err := uc.bankRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return banks, nil
}
