package servicerepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ozonmp/ssn-service-api/internal/metrics"

	"github.com/ozonmp/ssn-service-api/internal/tracer"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
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
	Remove(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) error
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
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Describe")
	defer sp.Finish()

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
		metrics.AddNotFoundErrorsTotal(1)
		return nil, ErrNoService
	}

	return &service, err
}

// Add - добавляет в репозиторий сервис.
func (r *serviceRepo) Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Add")
	defer sp.Finish()

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
		if rows == nil {
			return
		}

		if errCl := rows.Close(); errCl != nil {
			logger.ErrorKV(ctx, "serviceRepo.Add - failed close rows", "err", errCl)
		}
	}()

	if err != nil {
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&service.ID); err != nil {
			return err
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

// Update - обновляет сервис.
func (r *serviceRepo) Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Update")
	defer sp.Finish()

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
		metrics.AddNotFoundErrorsTotal(1)
		return ErrNoService
	}

	return nil
}

// List - возвращает постраничный список сервисов.
func (r *serviceRepo) List(ctx context.Context, offset uint64, limit uint64, tx repo.QueryerExecer) ([]*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.List")
	defer sp.Finish()

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

	return res, nil
}

// Remove - удаляет из репозитория сервис.
// Возвращает true если сервис существовал в репозитории и успешно удален методом.
func (r *serviceRepo) Remove(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) error {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Remove")
	defer sp.Finish()

	execer := r.getExecer(tx)

	query := sq.Update("services").PlaceholderFormat(sq.Dollar)
	query = query.Set("is_removed", true)
	query = query.Where(sq.And{sq.Eq{"id": serviceID}, sq.Eq{"is_removed": false}})

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
		metrics.AddNotFoundErrorsTotal(1)
		return ErrNoService
	}

	return nil
}

func (r *serviceRepo) getExecer(tx repo.QueryerExecer) repo.QueryerExecer {
	if tx != nil {
		return tx
	}

	return r.db
}
