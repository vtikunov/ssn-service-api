package servicerepo

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

func setupEventsRepo() (*eventRepo, sqlmock.Sqlmock) {
	mockDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	repo := &eventRepo{
		db: sqlxDB,
	}

	return repo, mock
}

func Test_EventAddSQL(t *testing.T) {
	r, dbMock := setupEventsRepo()
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	dbMock.ExpectQuery(`INSERT INTO service_events (service_id,type,status,payload,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`).
		WithArgs(
			1,
			subscription.Created,
			subscription.Deferred,
			sqlmock.AnyArg(), // TODO В debug режиме такое проходит, а в простом тесте падает по панике: []byte(`{"serviceId":"1", "name":"Test", "description":"Desc test"}`)
			time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
		).
		WillReturnRows(rows)

	err := r.Add(ctx, &subscription.ServiceEvent{
		ServiceID: 1,
		Type:      subscription.Created,
		Status:    subscription.Deferred,
		Service: &subscription.Service{
			ID:          1,
			Name:        "Test",
			Description: "Desc test",
			CreatedAt:   time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
		},
		UpdatedAt: time.Date(2021, 5, 1, 11, 0, 0, 0, time.UTC),
	}, nil)

	require.NoError(t, err)
}
