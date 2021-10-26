package repo

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// EventRepo - интерфейс репозитория событий
//
// Lock: блокирует и возвращает из репозитория n событий.
//
// Unlock: разблокирует в репозитории события с указанными eventIDs.
//
// Add: добавляет событие в репозиторий.
//
// Remove: удаляет из репозитория события с указанными eventIDs.
//
type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error

	Add(ctx context.Context, event []subscription.ServiceEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}
