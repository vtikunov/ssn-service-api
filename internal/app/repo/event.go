package repo

import (
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
	Lock(n uint64) ([]subscription.ServiceEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []subscription.ServiceEvent) error
	Remove(eventIDs []uint64) error
}
