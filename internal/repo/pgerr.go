package repo

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PostgresUniqueViolation     = "23505"
	PostgresForeignKeyViolation = "23503"
	PostgresCheckViolation      = "23514"
)

func IsPostgresErrorCode(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}

func IsUniqueViolation(err error) bool {
	return IsPostgresErrorCode(err, PostgresUniqueViolation)
}

func IsForeignKeyViolation(err error) bool {
	return IsPostgresErrorCode(err, PostgresForeignKeyViolation)
}

func IsCheckViolation(err error) bool {
	return IsPostgresErrorCode(err, PostgresCheckViolation)
}
