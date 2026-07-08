package payment

import (
	"context"
	"payment_gateway/internal/repo"
)

type Repositories struct {
	Payment        repo.Payment
	PaymentAttempt repo.PaymentAttempt
}

type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context, repos Repositories) error) error
}
