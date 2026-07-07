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

	Payment interface {
		Create(ctx context.Context, p *entity.Payment, bank_id int64) (*entity.Payment, error)
		GetById(ctx context.Context, id int64) (*entity.Payment, error)
		List(ctx context.Context) ([]*entity.Payment, error)
	}
)
