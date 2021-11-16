package servicerepo

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/tracer"

	"google.golang.org/protobuf/encoding/protojson"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/repo"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
)

// EventRepo - интерфейс репозитория событий.
//
// Add: добавляет событие в репозиторий.
//
type EventRepo interface {
	Add(ctx context.Context, event *subscription.ServiceEvent, tx repo.QueryerExecer) error
}

type eventRepo struct {
	db repo.QueryerExecer
}

// NewEventRepo создаёт инстанс репозитория.
func NewEventRepo(db repo.QueryerExecer) *eventRepo {
	return &eventRepo{db: db}
}

// Add - добавляет события в репозиторий.
func (r *eventRepo) Add(ctx context.Context, event *subscription.ServiceEvent, tx repo.QueryerExecer) error {
	sp := tracer.StartSpanFromContext(ctx, "eventRepo.Add")
	defer sp.Finish()

	execer := r.getExecer(tx)

	pbSrvPayload := &pb.ServiceEventPayload{}

	if event.Service != nil {
		pbSrvPayload.ServiceId = event.Service.ID
		pbSrvPayload.Name = event.Service.Name
		pbSrvPayload.Description = event.Service.Description
	}

	payload, err := protojson.Marshal(pbSrvPayload)

	if err != nil {
		return err
	}

	query := sq.Insert("service_events").PlaceholderFormat(sq.Dollar)
	query = query.Columns("service_id", "type", "status", "payload", "updated_at")
	query = query.Values(event.ServiceID, event.Type, event.Status, payload, event.UpdatedAt)

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = execer.ExecContext(ctx, s, args...)

	if err != nil {
		return err
	}

	return nil
}

func (r *eventRepo) getExecer(tx repo.QueryerExecer) repo.QueryerExecer {
	if tx != nil {
		return tx
	}

	return r.db
}
