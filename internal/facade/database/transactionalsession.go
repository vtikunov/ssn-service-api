package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/ssn-service-api/internal/facade/repo"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

type transactionalSession struct {
	db *sqlx.DB
}

// NewTransactionalSession - создаёт экзекьютора транзакций.
func NewTransactionalSession(db *sqlx.DB) *transactionalSession {
	return &transactionalSession{
		db: db,
	}
}

func (ts *transactionalSession) Execute(ctx context.Context, fn func(ctx context.Context, tx repo.QueryerExecer) error) error {
	tx, err := ts.db.BeginTxx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if err := fn(ctx, tx); err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			logger.ErrorKV(ctx, "rollback transaction failed", "err", errRb)
		}

		return err
	}

	return tx.Commit()
}
