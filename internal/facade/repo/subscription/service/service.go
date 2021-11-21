package servicerepo

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/facade/repo"
	"github.com/ozonmp/ssn-service-api/internal/tracer"
)

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

func (r *serviceRepo) getExecer(tx repo.QueryerExecer) repo.QueryerExecer {
	if tx != nil {
		return tx
	}

	return r.db
}
