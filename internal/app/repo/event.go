package repo

import (
	"context"
	"errors"
	"time"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/repo"
)

// EventRepo - интерфейс репозитория событий
//
// Lock: блокирует и возвращает из репозитория n событий.
//
// LockExceptLockedByServiceID: блокирует и возвращает из репозитория n событий.
// Метод учитывает блокировку по ID сервиса и блокирует только те события, ID сервисов которых
// не имеют ни одного заблокированного события.
//
// LockByServiceID: блокирует и возвращает из репозитория события по ID сервиса.
// Это отдельный вид блокировки - блокировка по ID события - не учитывается.
// Реализация метода должна учитывать, что если в репозитарии есть уже заблокированные события
// с переданным ID сервиса - блокировка незаблокированных невозможна. Также метод
// обязан гарантировать, что возвращаемые значения отсортированы в порядке их возникновения.
//
// Unlock: разблокирует в репозитории события с указанными eventIDs.
//
// Add: добавляет событие в репозиторий.
//
// Remove: удаляет из репозитория события с указанными eventIDs.
//
type EventRepo interface {
	Lock(ctx context.Context, n uint64, tx QueryerExecer) ([]subscription.ServiceEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64, tx QueryerExecer) error

	Remove(ctx context.Context, eventIDs []uint64, tx QueryerExecer) error
}

type QueryerExecer interface {
	sqlx.Execer
	sqlx.Queryer
	sqlx.QueryerContext
	sqlx.ExecerContext
}

// ErrNoEvent - ошибка отсутствия события в репозитории.
var ErrNoEvent = errors.New("event is not exists")

type eventRepo struct {
	db QueryerExecer
}

// NewEventRepo создаёт инстанс репозитория.
func NewEventRepo(db repo.QueryerExecer) *eventRepo {
	return &eventRepo{db: db}
}

func (r *eventRepo) getExecer(tx repo.QueryerExecer) repo.QueryerExecer {
	if tx != nil {
		return tx
	}

	return r.db
}

type serviceEvent struct {
	ID        uint64                   `db:"id"`
	ServiceID uint64                   `db:"service_id"`
	Type      subscription.EventType   `db:"type"`
	Status    subscription.EventStatus `db:"status"`
	Service   servicePayload           `db:"payload"`
	UpdatedAt time.Time                `db:"updated_at"`
}

type servicePayload pb.ServiceEventPayload

func (sp *servicePayload) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("incompatible type for servicePayload")
	}

	res := &pb.ServiceEventPayload{}

	err := protojson.Unmarshal(source, res)

	if err != nil {
		return err
	}

	*sp = servicePayload{
		ServiceId:   res.ServiceId,
		Name:        res.Name,
		Description: res.Description,
	}

	return nil
}

func (se *serviceEvent) convertToServiceEvent() *subscription.ServiceEvent {
	return &subscription.ServiceEvent{
		ID:        se.ID,
		ServiceID: se.ServiceID,
		Type:      se.Type,
		Status:    se.Status,
		Service: &subscription.Service{
			ID:          se.Service.ServiceId,
			Name:        se.Service.Name,
			Description: se.Service.Description,
		},
		UpdatedAt: se.UpdatedAt,
	}
}
