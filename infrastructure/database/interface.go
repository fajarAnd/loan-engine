package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	Querier interface {
		// QueryRow executes a query that is expected to return at most one row.
		Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
		// QueryRow returns a single row from the database. If no rows are returned, it will return a row with all zero values.
		QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	}

	Executor interface {
		Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	}

	Tx interface {
		Begin(ctx context.Context) (pgx.Tx, error)
		Querier
	}
)
