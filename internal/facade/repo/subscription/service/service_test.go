package servicerepo

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
)

func setupServiceRepo() (*serviceRepo, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repo := NewServiceRepo(sqlxDB)

	return repo, mock
}

func Test_ServiceAddSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	dbMock.ExpectExec(`INSERT INTO services (id,name,description,last_event_id) VALUES ($1,$2,$3,$4) ON CONFLICT (id) DO NOTHING`).
		WithArgs(
			uint64(1),
			"Test",
			"Test Desc",
			uint64(2),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := r.Add(ctx, &subscription.Service{
		ID:          uint64(1),
		Name:        "Test",
		Description: "Test Desc",
		LastEventID: uint64(2),
	}, nil)

	require.NoError(t, err)
}

func Test_ServiceUpdateSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	dbMock.ExpectExec(`UPDATE services SET name = $1, description = $2, last_event_id = $3 WHERE (id = $4 AND is_removed = $5 AND last_event_id < $6)`).
		WithArgs(
			"Test",
			"Test Desc",
			uint64(3),
			uint64(1),
			false,
			uint64(3),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := r.Update(ctx, &subscription.Service{
		ID:          uint64(1),
		Name:        "Test",
		Description: "Test Desc",
		LastEventID: uint64(3),
	}, nil)

	require.NoError(t, err)
}

func Test_ServiceRemoveSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	dbMock.ExpectExec(`UPDATE services SET is_removed = $1, last_event_id = $2 WHERE (id = $3 AND is_removed = $4)`).
		WithArgs(true, uint64(4), uint64(1), false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := r.Remove(ctx, uint64(1), uint64(4), nil)

	require.NoError(t, err)
}
