package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Exec(
		ctx context.Context,
		sql string,
		arguments ...interface{},
	) (pgconn.CommandTag, error)

	QueryRow(
		ctx context.Context,
		sql string,
		args ...interface{},
	) pgx.Row
	Query(
		ctx context.Context,
		sql string,
		args ...interface{},
	) (pgx.Rows, error)
}
