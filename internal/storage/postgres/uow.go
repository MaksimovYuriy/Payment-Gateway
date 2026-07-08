package postgres

import (
	"context"
	"database/sql"
	prepo "payment_gateway/internal/repo/payment"
	parepo "payment_gateway/internal/repo/payment_attempt"
	"payment_gateway/internal/usecase/payment"
)

type UnitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos payment.Repositories) error) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	repos := payment.Repositories{
		Payment:        prepo.NewRepo(tx),
		PaymentAttempt: parepo.NewRepo(tx),
	}

	if err := fn(ctx, repos); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
