package merchant

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

var _ repo.Merchant = (*Repo)(nil)

func (r *Repo) Create(ctx context.Context, m *entity.Merchant) error {
	query := `
		INSERT INTO merchants (name, domain, api_key_hash,
							   webhook_url, success_redirect_url, 
							   failure_redirect_url, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	row := r.database.QueryRowContext(
		ctx, query, m.Name, m.Domain, m.ApiKeyHash,
		m.WebhookUrl, m.SuccessRedirectUrl, m.FailureRedirectUrl,
		m.IsActive,
	)

	if err := row.Scan(&m.Id, &m.CreatedAt, &m.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (r *Repo) List(ctx context.Context) ([]*entity.Merchant, error) {
	query := `
		SELECT id, name, domain, webhook_url, 
			   success_redirect_url, 
			   failure_redirect_url, is_active, 
			   created_at, updated_at
		FROM merchants
	`
	rows, err := r.database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	merchants := make([]*entity.Merchant, 0)
	for rows.Next() {
		var merchant entity.Merchant
		err := rows.Scan(
			&merchant.Id,
			&merchant.Name,
			&merchant.Domain,
			&merchant.WebhookUrl,
			&merchant.SuccessRedirectUrl,
			&merchant.FailureRedirectUrl,
			&merchant.IsActive,
			&merchant.CreatedAt,
			&merchant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		merchants = append(merchants, &merchant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return merchants, nil
}
