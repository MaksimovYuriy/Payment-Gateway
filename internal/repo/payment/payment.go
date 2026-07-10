package payment

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
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Repo struct {
	database executor
}

func NewRepo(db executor) *Repo {
	return &Repo{database: db}
}

var _ repo.Payment = (*Repo)(nil)

func (r *Repo) Create(ctx context.Context, p *entity.Payment) error {
	query := `
		INSERT INTO payments (
								merchant_id, order_id, amount,
								currency
							 )
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, created_at, updated_at
	`
	row := r.database.QueryRowContext(ctx, query, p.MerchantId, p.OrderId, p.Amount,
		p.Currency)

	if err := row.Scan(&p.Id, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if repo.IsUniqueViolation(err) {
			return apperr.AlreadyExists("payment already exists")
		}
		if repo.IsForeignKeyViolation(err) {
			return apperr.InvalidInput("merchant not found")
		}
		if repo.IsCheckViolation(err) {
			return apperr.InvalidInput(apperr.MessageInvalidInput)
		}
		return apperr.Internal("failed to create payment", err)
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, p *entity.Payment) error {
	query := `
		UPDATE payments
		SET success_bank_id = $1, status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING updated_at
	`

	row := r.database.QueryRowContext(ctx, query, p.SuccessBankId, p.Status, p.Id)
	if err := row.Scan(&p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.NotFound("payment not found")
		}
		if repo.IsForeignKeyViolation(err) {
			return apperr.InvalidInput("bank not found")
		}
		if repo.IsCheckViolation(err) {
			return apperr.InvalidInput(apperr.MessageInvalidInput)
		}
		return apperr.Internal("failed to update payment", err)
	}
	return nil
}

func (r *Repo) GetById(ctx context.Context, id int64) (*entity.Payment, error) {
	query := `
		SELECT id, merchant_id, success_bank_id, order_id,
			   amount, currency, status, created_at, updated_at
		FROM payments
		WHERE id = $1
	`
	row := r.database.QueryRowContext(ctx, query, id)
	var payment entity.Payment
	if err := row.Scan(
		&payment.Id,
		&payment.MerchantId,
		&payment.SuccessBankId,
		&payment.OrderId,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperr.NotFound("payment not found")
		}
		return nil, apperr.Internal("failed to get payment", err)
	}
	return &payment, nil
}

func (r *Repo) List(ctx context.Context) ([]*entity.Payment, error) {
	query := `
		SELECT id, merchant_id, success_bank_id, order_id,
			   amount, currency, status, created_at, updated_at
		FROM payments
	`
	rows, err := r.database.QueryContext(ctx, query)
	if err != nil {
		return nil, apperr.Internal("failed to list payments", err)
	}
	defer rows.Close()
	payments := make([]*entity.Payment, 0)
	for rows.Next() {
		var payment entity.Payment
		if err := rows.Scan(
			&payment.Id,
			&payment.MerchantId,
			&payment.SuccessBankId,
			&payment.OrderId,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {
			return nil, apperr.Internal("failed to scan payment", err)
		}
		payments = append(payments, &payment)
	}
	if err := rows.Err(); err != nil {
		return nil, apperr.Internal("failed to iterate payments", err)
	}
	return payments, nil
}
