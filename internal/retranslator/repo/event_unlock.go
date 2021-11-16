package repo

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/retranslator/metrics"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

func (r *eventRepo) Unlock(ctx context.Context, eventIDs []uint64, tx QueryerExecer) error {
	execer := r.getExecer(tx)

	query := sq.Update("service_events").PlaceholderFormat(sq.Dollar)
	query = query.Set("status", subscription.Deferred)
	query = query.Set("updated_at", "NOW()")
	query = query.Where(sq.Eq{"id": eventIDs})

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
		return ErrNoEvent
	}

	metrics.SubEventsCountInPool(uint(len(eventIDs)))

	return nil
}
