package bank

import (
	"context"
	"database/sql"
	"payment_gateway/internal/entity"
	"payment_gateway/internal/repo"
)

type Repo struct {
	database *sql.DB
}

var _ repo.Bank = (*Repo)(nil)

func NewRepo(db *sql.DB) *Repo {
	return &Repo{database: db}
}

func (r *Repo) Create(ctx context.Context, b *entity.Bank) error {
	const query = `
		INSERT INTO banks (code, name, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.database.QueryRowContext(
		ctx, query, b.Code, b.Name, b.IsActive,
	).Scan(&b.Id, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) List(ctx context.Context) ([]*entity.Bank, error) {
	const query = `
		SELECT id, code, name, is_active, created_at, updated_at
		FROM banks
	`
	rows, err := r.database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banks := make([]*entity.Bank, 0)
	for rows.Next() {
		var bank entity.Bank
		if err := rows.Scan(
			&bank.Id,
			&bank.Code,
			&bank.Name,
			&bank.IsActive,
			&bank.CreatedAt,
			&bank.UpdatedAt,
		); err != nil {
			return nil, err
		}
		banks = append(banks, &bank)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return banks, nil
}
