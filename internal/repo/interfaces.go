package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// QueryerExecer - универсальный интерфейс для обычного соединения и соединения-транзакции.
type QueryerExecer interface {
	sqlx.Execer
	sqlx.Queryer
	sqlx.QueryerContext
	sqlx.ExecerContext
}

// TransactionalSession - оборачивает в транзакцию функцию fn.
type TransactionalSession interface {
	Execute(ctx context.Context, fn func(ctx context.Context, tx QueryerExecer) error) error
}
