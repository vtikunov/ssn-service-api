package repo

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// EventRepo - интерфейс репозитория событий
//
// Lock: блокирует и возвращает из репозитория n событий.
//
// LockByServiceID: блокирует и возвращает из репозитория события по ID сервиса.
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
	Lock(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error)
	LockByServiceID(ctx context.Context, serviceID uint64) ([]subscription.ServiceEvent, error)
	Unlock(eventIDs []uint64) error

	Add(ctx context.Context, event []subscription.ServiceEvent) error
	Remove(eventIDs []uint64) error
}
