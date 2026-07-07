package payment

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

var _ repo.Payment = (*Repo)(nil)

func (r *Repo) Create(ctx context.Context, p *entity.Payment) error {
	query := `
		INSERT INTO payments (
								merchant_id, order_id, amount,
								currency, created_at, updated_at 
							 )
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, status
	`
	row := r.database.QueryRowContext(ctx, query, p.MerchantId, p.OrderId, p.Amount,
		p.Currency, p.CreatedAt, p.UpdatedAt)

	if err := row.Scan(&p.Id, &p.Status); err != nil {
		return err
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
		return err
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
		return nil, err
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
		return nil, err
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
			return nil, err
		}
		payments = append(payments, &payment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return payments, nil
}
