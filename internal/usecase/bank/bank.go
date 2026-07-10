package bank

import (
	"context"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/repo"
	"payment_gateway/internal/usecase"
)

type UseCase struct {
	bankRepo repo.Bank
}

var _ usecase.Bank = (*UseCase)(nil)

func NewUseCase(br repo.Bank) *UseCase {
	return &UseCase{
		bankRepo: br,
	}
}

func (uc *UseCase) Registration(ctx context.Context, br *entity.Bank) (*entity.Bank, error) {
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
