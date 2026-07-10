package paymentattempt

import (
	"context"
	"database/sql"
	"errors"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/lib/apperr"
	"payment_gateway/internal/repo"
)

type executor interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Repo struct {
	database executor
}

func NewRepo(db executor) *Repo {
	return &Repo{database: db}
}

var _ repo.PaymentAttempt = (*Repo)(nil)

func (r *Repo) Create(ctx context.Context, pa *entity.PaymentAttempt) error {
	query := `
		INSERT INTO payment_attempts (payment_id, bank_id)
		VALUES ($1, $2)
		RETURNING id, status, created_at, updated_at
	`
	row := r.database.QueryRowContext(ctx, query, pa.PaymentId, pa.BankId)
	if err := row.Scan(&pa.Id, &pa.Status, &pa.CreatedAt, &pa.UpdatedAt); err != nil {
		if repo.IsForeignKeyViolation(err) {
			return apperr.InvalidInput("payment or bank not found")
		}
		if repo.IsCheckViolation(err) {
			return apperr.InvalidInput(apperr.MessageInvalidInput)
		}
		return apperr.Internal("failed to create payment attempt", err)
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, pa *entity.PaymentAttempt) error {
	query := `
		UPDATE payment_attempts
		SET status = $1, external_payment_id = $2, error_message = $3, error_code = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`
	row := r.database.QueryRowContext(ctx, query, pa.Status, pa.ExternalPaymentId, pa.ErrorMessage, pa.ErrorCode, pa.Id)
	if err := row.Scan(&pa.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.NotFound("payment attempt not found")
		}
		if repo.IsCheckViolation(err) {
			return apperr.InvalidInput(apperr.MessageInvalidInput)
		}
		return apperr.Internal("failed to update payment attempt", err)
	}
	return nil
}
