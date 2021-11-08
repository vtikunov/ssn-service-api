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

func setup() (*eventRepo, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repo := &eventRepo{
		db: sqlxDB,
	}

	return repo, mock
}

func Test_LockSQL(t *testing.T) {
	r, dbMock := setup()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "service_id", "type", "status", "payload", "updated_at"}).
		AddRow(1, 2, subscription.Created, subscription.Processed, "{}", time.Now())

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
				RETURNING s.id, s.service_id, s.type, s.status, s.payload, s.updated_at`).
		WithArgs(subscription.Processed, "NOW()", subscription.Processed, subscription.Deferred).
		WillReturnRows(rows)

	_, err := r.Lock(ctx, 2, nil)

	require.NoError(t, err)
}
