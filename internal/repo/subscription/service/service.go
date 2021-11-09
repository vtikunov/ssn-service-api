package servicerepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/repo"
)

// ErrNoService - ошибка отсутствия сервиса в репозитории.
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
	Describe(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (*subscription.Service, error)
	Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error
	Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error
	List(ctx context.Context, offset uint64, limit uint64, tx repo.QueryerExecer) ([]*subscription.Service, error)
	Remove(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (ok bool, err error)
}

type serviceRepo struct {
	db repo.QueryerExecer
}

// NewServiceRepo создаёт инстанс репозитория.
func NewServiceRepo(db repo.QueryerExecer) *serviceRepo {
	return &serviceRepo{db: db}
}

// Describe - возвращает из репозитория сервис по его ID.
func (r *serviceRepo) Describe(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (*subscription.Service, error) {
	execer := r.getExecer(tx)

	query := sq.Select("id, name, description, is_removed, created_at, updated_at").PlaceholderFormat(sq.Dollar).From("services")
	query = query.Where(sq.And{sq.Eq{"id": serviceID}, sq.Eq{"is_removed": false}})

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var service subscription.Service
	err = execer.QueryRowxContext(ctx, s, args...).StructScan(&service)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoService
	} else if err != nil {
		return nil, err
	}

	return &service, err
}

// Add - добавляет в репозиторий сервис.
func (r *serviceRepo) Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error {
	execer := r.getExecer(tx)

	query := sq.Insert("services").PlaceholderFormat(sq.Dollar)
	query = query.Columns("name", "description", "created_at", "updated_at")
	query = query.Values(service.Name, service.Description, service.CreatedAt, service.UpdatedAt)
	query = query.Suffix("RETURNING id")

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	rows, err := execer.QueryContext(ctx, s, args...)
	defer func() {
		if errCl := rows.Close(); errCl != nil {
			log.Error().Err(err)
		}
	}()

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
func (r *serviceRepo) Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error {
	execer := r.getExecer(tx)

	query := sq.Update("services").PlaceholderFormat(sq.Dollar)
	query = query.Set("name", service.Name)
	query = query.Set("description", service.Description)
	query = query.Set("updated_at", service.UpdatedAt)
	query = query.Where(sq.And{sq.Eq{"id": service.ID}, sq.Eq{"is_removed": false}})

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	res, err := execer.ExecContext(ctx, s, args...)

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
func (r *serviceRepo) List(ctx context.Context, offset uint64, limit uint64, tx repo.QueryerExecer) ([]*subscription.Service, error) {
	execer := r.getExecer(tx)

	query := sq.Select("*").PlaceholderFormat(sq.Dollar).From("services")
	query = query.Where(sq.Eq{"is_removed": false})
	query = query.OrderBy("id ASC")
	query = query.Offset(offset)
	query = query.Limit(limit)

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := execer.QueryContext(ctx, s, args...)

	if err != nil {
		return nil, err
	}

	res := make([]*subscription.Service, 0)
	err = sqlx.StructScan(rows, &res)

	if err != nil {
		return nil, err
	}

	return res, err
}

// Remove - удаляет из репозитория сервис.
// Возвращает true если сервис существовал в репозитории и успешно удален методом.
func (r *serviceRepo) Remove(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (ok bool, err error) {
	execer := r.getExecer(tx)

	query := sq.Update("services").PlaceholderFormat(sq.Dollar)
	query = query.Set("is_removed", true)
	query = query.Where(sq.And{sq.Eq{"id": serviceID}, sq.Eq{"is_removed": false}})

	s, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	res, err := execer.ExecContext(ctx, s, args...)

	if err != nil {
		return false, err
	}

	num, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	return num > 0, nil
}

func (r *serviceRepo) getExecer(tx repo.QueryerExecer) repo.QueryerExecer {
	if tx != nil {
		return tx
	}

	return r.db
}
