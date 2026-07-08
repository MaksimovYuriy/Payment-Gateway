package merchant

import (
	"context"
	"payment_gateway/internal/entity"
	apikey "payment_gateway/internal/lib/api_key"
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

func (uc *UseCase) Registration(ctx context.Context, m *entity.Merchant) (*entity.Merchant, string, error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	apiKey, err := apikey.GenerateApiKey()
	if err != nil {
		return nil, "", err
	}
	m.ApiKeyHash, err = apikey.HashApiKey(apiKey)
	if err != nil {
		return nil, "", err
	}

	if err := uc.merchantRepo.Create(ctx, m); err != nil {
		return nil, "", err
	}
	return m, apiKey, nil
}

func (uc *UseCase) List(ctx context.Context) ([]*entity.Merchant, error) {
	merchants, err := uc.merchantRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	return merchants, nil
}
