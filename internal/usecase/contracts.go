package usecase

import (
	"context"
	"payment_gateway/internal/entity"
)

type (
	Bank interface {
		Registration(ctx context.Context, b *entity.Bank) (*entity.Bank, error)
		List(ctx context.Context) ([]*entity.Bank, error)
	}

	Merchant interface {
		Registration(ctx context.Context, m *entity.Merchant) (*entity.Merchant, error)
		List(ctx context.Context) ([]*entity.Merchant, error)
	}
)
