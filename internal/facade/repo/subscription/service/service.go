package servicerepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/facade/repo"
	"github.com/ozonmp/ssn-service-api/internal/tracer"
)

// ErrNoService - ошибка отсутствия сервиса в репозитории.
var ErrNoService = errors.New("service is not exists")

// ServiceRepo - Репозиторий сервисов
//
// Add: добавляет в репозиторий сервис.
//
// Update: обновляет сервис.
//
// Remove: удаляет из репозитория сервис.
//
type ServiceRepo interface {
	Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error)
	Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error)
	Remove(ctx context.Context, serviceID, eventID uint64, tx repo.QueryerExecer) (bool, error)
}

type serviceRepo struct {
	db repo.QueryerExecer
}

// NewServiceRepo создаёт инстанс репозитория.
func NewServiceRepo(db repo.QueryerExecer) *serviceRepo {
	return &serviceRepo{db: db}
}

// Add - добавляет в репозиторий сервис.
func (r *serviceRepo) Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error) {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Add")
	defer sp.Finish()

	execer := r.getExecer(tx)

	query := sq.Insert("services").PlaceholderFormat(sq.Dollar)
	query = query.Columns("id", "name", "description", "last_event_id")
	query = query.Values(service.ID, service.Name, service.Description, service.LastEventID)
	query = query.Suffix("ON CONFLICT (id) DO NOTHING")

	s, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	res, err := execer.ExecContext(ctx, s, args...)

	if err != nil {
		return false, err
	}

	num, err := res.RowsAffected()

	return num > 0, err
}

// Update - обновляет сервис.
func (r *serviceRepo) Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error) {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Update")
	defer sp.Finish()

	execer := r.getExecer(tx)

	query := sq.Update("services").PlaceholderFormat(sq.Dollar)
	query = query.Set("name", service.Name)
	query = query.Set("description", service.Description)
	query = query.Set("last_event_id", service.LastEventID)
	query = query.Where(sq.And{sq.Eq{"id": service.ID}, sq.Eq{"is_removed": false}, sq.Lt{"last_event_id": service.LastEventID}})

	s, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	res, err := execer.ExecContext(ctx, s, args...)

	if err != nil {
		return false, err
	}

	num, err := res.RowsAffected()

	return num > 0, err
}

// Remove - удаляет из репозитория сервис.
func (r *serviceRepo) Remove(ctx context.Context, serviceID, eventID uint64, tx repo.QueryerExecer) (bool, error) {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Remove")
	defer sp.Finish()

	execer := r.getExecer(tx)

	query := sq.Update("services").PlaceholderFormat(sq.Dollar)
	query = query.Set("is_removed", true)
	query = query.Set("last_event_id", eventID)
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

	return num > 0, err
}

func (r *serviceRepo) Describe(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "serviceRepo.Describe")
	defer sp.Finish()

	execer := r.getExecer(tx)

	query := sq.Select("*").PlaceholderFormat(sq.Dollar).From("services")
	query = query.Where(sq.And{sq.Eq{"id": serviceID}, sq.Eq{"is_removed": false}})

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var service subscription.Service
	err = execer.QueryRowxContext(ctx, s, args...).StructScan(&service)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoService
	}

	return &service, err
}

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

func (r *serviceRepo) getExecer(tx repo.QueryerExecer) repo.QueryerExecer {
	if tx != nil {
		return tx
	}

	return r.db
}
