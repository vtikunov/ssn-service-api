package repo

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

type EventRepo interface {
	Add(ctx context.Context, event []subscription.ServiceEvent, tx QueryerExecer) error
}

type eventRepo struct {
	db QueryerExecer
}

// NewEventRepo создаёт инстанс репозитория.
func NewEventRepo(db QueryerExecer) *eventRepo {
	return &eventRepo{db: db}
}

// Add - добавляет события в репозиторий.
func (r *eventRepo) Add(ctx context.Context, event []subscription.ServiceEvent, tx QueryerExecer) error {
	panic("implement me")

	return ErrNoService
}

func (r *eventRepo) getExecer(tx QueryerExecer) QueryerExecer {
	if tx != nil {
		return tx
	}

	return r.db
}
