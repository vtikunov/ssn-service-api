package repo

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

func setupEventRepo() (*eventRepo, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repo := &eventRepo{
		db: sqlxDB,
	}

	return repo, mock
}

func Test_EventsLockSQL(t *testing.T) {
	r, dbMock := setupEventRepo()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "service_id", "type", "status", "payload", "updated_at"}).
		AddRow(1, 2, subscription.Created, subscription.Processed, "{}", time.Now())

	dbMock.ExpectExec("select pg_try_advisory_xact_lock(1936028278, 1819239275)").
		WillReturnResult(sqlmock.NewResult(0, 0))

	dbMock.ExpectQuery(`
				UPDATE service_events
				SET status = $1,
					updated_at = $2
				WHERE id IN (
					SELECT s.id
					FROM service_events s
						 LEFT JOIN service_events s1 on (s1.service_id = s.service_id AND s1.status = $3)
					WHERE (s1.service_id IS NULL AND s.status = $4)
					ORDER BY s.id
					LIMIT 2 )
				AND status = $5
				RETURNING id, service_id, type, status, payload, updated_at`).
		WithArgs(subscription.Processed, "NOW()", subscription.Processed, subscription.Deferred, subscription.Deferred).
		WillReturnRows(rows)

	_, err := r.Lock(ctx, 2, nil)

	require.NoError(t, err)
}

func Test_EventsUnlockSQL(t *testing.T) {
	r, dbMock := setupEventRepo()
	ctx := context.Background()

	dbMock.ExpectExec(`UPDATE service_events SET status = $1, updated_at = $2 WHERE id IN ($3,$4)`).
		WithArgs(subscription.Deferred, "NOW()", 1, 2).
		WillReturnResult(sqlmock.NewResult(2, 2))

	err := r.Unlock(ctx, []uint64{1, 2}, nil)

	require.NoError(t, err)
}

func Test_EventsRemoveSQL(t *testing.T) {
	r, dbMock := setupEventRepo()
	ctx := context.Background()

	dbMock.ExpectExec(`DELETE FROM service_events WHERE id IN ($1,$2)`).
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(2, 2))

	err := r.Remove(ctx, []uint64{1, 2}, nil)

	require.NoError(t, err)
}
