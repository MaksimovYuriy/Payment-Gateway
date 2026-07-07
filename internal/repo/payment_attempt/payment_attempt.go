package paymentattempt

import (
	"context"
	"database/sql"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/repo"
)

type Repo struct {
	database *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
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
		return err
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
		return err
	}
	return nil
}
