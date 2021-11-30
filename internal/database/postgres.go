package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

// NewPostgres returns DB
func NewPostgres(ctx context.Context, dsn, driver string, retryAttempts uint64) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		logger.ErrorKV(ctx, "failed to create database connection", "err", err)

		return nil, err
	}

	var attempt uint64 = 1

	if err = db.PingContext(ctx); err != nil {

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if attempt >= retryAttempts {
				logger.ErrorKV(ctx, "failed ping the database", "err", err, "attempts", attempt)
				return nil, err
			}
			attempt++

			logger.InfoKV(ctx, "reconnecting to DB...")

			if err = db.PingContext(ctx); err == nil {
				return db, nil
			}

			logger.InfoKV(ctx, "connection was lost. waiting for 1 sec...")
		}
	}

	return db, nil
}
