package servicerepo

import (
	"context"
	"testing"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupServiceRepo() (*serviceRepo, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repo := &serviceRepo{
		db: sqlxDB,
	}

	return repo, mock
}

func Test_ServiceDescribeSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "is_removed", "created_at", "updated_at"}).
		AddRow(1, "Test", "Test desc", false, time.Now(), time.Now())

	dbMock.ExpectQuery(`SELECT id, name, description, is_removed, created_at, updated_at FROM services WHERE (id = $1 AND is_removed = $2)`).
		WithArgs(1, false).
		WillReturnRows(rows)

	_, err := r.Describe(ctx, 1, nil)

	require.NoError(t, err)
}

func Test_ServiceAddSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	dbMock.ExpectQuery(`INSERT INTO services (name,description,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING id`).
		WithArgs(
			"Test",
			"Test Desc",
			time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
			time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
		).
		WillReturnRows(rows)

	err := r.Add(ctx, &subscription.Service{
		Name:        "Test",
		Description: "Test Desc",
		CreatedAt:   time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
	}, nil)

	require.NoError(t, err)
}

func Test_ServiceUpdateSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	dbMock.ExpectExec(`UPDATE services SET name = $1, description = $2, updated_at = $3 WHERE (id = $4 AND is_removed = $5)`).
		WithArgs(
			"Test",
			"Test Desc",
			time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
			1,
			false,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.Update(ctx, &subscription.Service{
		ID:          1,
		Name:        "Test",
		Description: "Test Desc",
		UpdatedAt:   time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
	}, nil)

	require.NoError(t, err)
}

func Test_ServiceListSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "description", "is_removed", "created_at", "updated_at"}).
		AddRow(1, "Test", "Test desc", false, time.Now(), time.Now()).
		AddRow(2, "Test2", "Test desc2", false, time.Now(), time.Now())

	dbMock.ExpectQuery(`SELECT * FROM services WHERE is_removed = $1 ORDER BY id ASC LIMIT 2 OFFSET 0`).
		WithArgs(false).
		WillReturnRows(rows)

	_, err := r.List(ctx, 0, 2, nil)

	require.NoError(t, err)
}

func Test_ServiceRemoveSQL(t *testing.T) {
	r, dbMock := setupServiceRepo()
	ctx := context.Background()

	dbMock.ExpectExec(`UPDATE services SET is_removed = $1 WHERE (id = $2 AND is_removed = $3)`).
		WithArgs(true, 1, false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := r.Remove(ctx, 1, nil)

	require.NoError(t, err)
}
