package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var ErrNoService = errors.New("service is not exists")

// ServiceRepo - Репозиторий сервисов
//
// Describe: возвращает из репозитория сервис по его ID.
//
// Add: добавляет в репозиторий сервис.
//
// Update: обновляет сервис.
//
// List: возвращает постраничный список сервисов.
//
// Remove: удаляет из репозитория сервис.
// Возвращает true если сервис существовал в репозитории и успешно удален методом.
//
type ServiceRepo interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	Add(ctx context.Context, service *subscription.Service) error
	Update(ctx context.Context, service *subscription.Service) error
	List(ctx context.Context) ([]*subscription.Service, error)
	Remove(ctx context.Context, serviceID uint64) (ok bool, err error)
}

type repo struct {
	db *sqlx.DB
}

// NewRepo создаёт инстанс репозитория
func NewRepo(db *sqlx.DB) *repo {
	return &repo{db: db}
}

// Describe - возвращает из репозитория сервис по его ID.
func (r repo) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	query := sq.Select("*").PlaceholderFormat(sq.Dollar).From("services")
	query = query.Where(sq.And{sq.Eq{"id": serviceID}, sq.Eq{"is_removed": false}})

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var service subscription.Service
	err = r.db.QueryRowxContext(ctx, s, args...).StructScan(&service)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoService
	} else if err != nil {
		return nil, err
	}

	return &service, err
}

// Add - добавляет в репозиторий сервис.
func (r repo) Add(ctx context.Context, service *subscription.Service) error {
	query := sq.Insert("services").PlaceholderFormat(sq.Dollar)
	query = query.Columns("name", "description", "created_at", "updated_at")
	query = query.Values(service.Name, service.Description, service.CreatedAt, service.UpdatedAt)
	query = query.Suffix("RETURNING id").RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(&service.ID)

		if err != nil {
			return err
		}

		return nil
	}

	return ErrNoService
}

// Update - обновляет сервис.
func (r repo) Update(ctx context.Context, service *subscription.Service) error {
	query := sq.Update("services").PlaceholderFormat(sq.Dollar)
	query = query.Set("name", service.Name)
	query = query.Set("description", service.Description)
	query = query.Set("updated_at", service.UpdatedAt)
	query = query.Where(sq.And{sq.Eq{"id": service.ID}, sq.Eq{"is_removed": false}})

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, s, args...)

	if err != nil {
		return err
	}

	num, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if num == 0 {
		return ErrNoService
	}

	return nil
}

// List - возвращает постраничный список сервисов.
func (r repo) List(ctx context.Context) ([]*subscription.Service, error) {
	query := sq.Select("*").PlaceholderFormat(sq.Dollar).From("services")
	query = query.Where(sq.Eq{"is_removed": false})
	query = query.OrderBy("id ASC")

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	res := make([]*subscription.Service, 0)
	err = r.db.SelectContext(ctx, &res, s, args...)

	return res, err
}

// Remove - удаляет из репозитория сервис.
// Возвращает true если сервис существовал в репозитории и успешно удален методом.
func (r repo) Remove(ctx context.Context, serviceID uint64) (ok bool, err error) {
	query := sq.Update("services").PlaceholderFormat(sq.Dollar)
	query = query.Set("is_removed", true)
	query = query.Where(sq.And{sq.Eq{"id": serviceID}, sq.Eq{"is_removed": false}})

	s, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	res, err := r.db.ExecContext(ctx, s, args...)

	if err != nil {
		return false, err
	}

	num, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	return num > 0, nil
}
