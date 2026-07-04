package repo

import (
	"context"
	"payment_gateway/internal/entity"
)

type (
	BankRepo interface {
		Create(ctx context.Context, b *entity.Bank) error
		List(ctx context.Context) ([]*entity.Bank, error)
	}
)
