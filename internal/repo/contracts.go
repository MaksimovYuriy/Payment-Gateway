package repo

import (
	"context"
	"payment_gateway/internal/entity"
)

type (
	Bank interface {
		Create(ctx context.Context, b *entity.Bank) error
		List(ctx context.Context) ([]*entity.Bank, error)
	}

	Merchant interface {
		Create(ctx context.Context, m *entity.Merchant) error
		List(ctx context.Context) ([]*entity.Merchant, error)
	}

	Payment interface {
		Create(ctx context.Context, p *entity.Payment) error
		Update(ctx context.Context, p *entity.Payment) error
		GetById(ctx context.Context, id int64) (*entity.Payment, error)
		List(ctx context.Context) ([]*entity.Payment, error)
	}

	PaymentAttempt interface {
		Create(ctx context.Context, pa *entity.PaymentAttempt) error
		Update(ctx context.Context, pa *entity.PaymentAttempt) error
	}
)
