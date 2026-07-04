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
)
