package merchant

import (
	"context"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apikey"
	"payment_gateway/internal/lib/apperr"
	"payment_gateway/internal/repo"
	"payment_gateway/internal/usecase"

	"github.com/go-playground/validator/v10"
)

type UseCase struct {
	merchantRepo repo.Merchant
	validate     *validator.Validate
}

func NewUseCase(repo repo.Merchant, vd *validator.Validate) *UseCase {
	return &UseCase{merchantRepo: repo, validate: vd}
}

var _ usecase.Merchant = (*UseCase)(nil)

func (uc *UseCase) Registration(ctx context.Context, m *entity.Merchant) (*entity.Merchant, string, error) {
	apiKey, err := apikey.GenerateApiKey()
	if err != nil {
		return nil, "", apperr.Internal("failed to generate api key", err)
	}
	m.ApiKeyHash, err = apikey.HashApiKey(apiKey)
	if err != nil {
		return nil, "", apperr.Internal("failed to hash api key", err)
	}
	if err := uc.validate.Struct(m); err != nil {
		return nil, "", apperr.InvalidInput(apperr.MessageInvalidInput)
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
