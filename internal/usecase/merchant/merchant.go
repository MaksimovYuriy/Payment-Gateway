package merchant

import (
	"context"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/repo"
	"payment_gateway/internal/usecase"
	"time"
)

type UseCase struct {
	merchantRepo repo.Merchant
}

func NewUseCase(repo repo.Merchant) *UseCase {
	return &UseCase{merchantRepo: repo}
}

var _ usecase.Merchant = (*UseCase)(nil)

func (uc *UseCase) Registration(ctx context.Context, m *entity.Merchant) (*entity.Merchant, error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	m.ApiKeyHash = "###" + m.Domain // Some service for generation api keys

	if err := uc.merchantRepo.Create(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (uc *UseCase) List(ctx context.Context) ([]*entity.Merchant, error) {
	merchants, err := uc.merchantRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return merchants, nil
}
