package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type QueryerExecer interface {
	sqlx.Execer
	sqlx.Queryer
	sqlx.QueryerContext
	sqlx.ExecerContext
}

type TransactionalSession interface {
	Execute(ctx context.Context, fn func(ctx context.Context, tx QueryerExecer) error) error
}
