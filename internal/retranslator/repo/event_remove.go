package repo

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/retranslator/metrics"

	sq "github.com/Masterminds/squirrel"
)

func (r *eventRepo) Remove(ctx context.Context, eventIDs []uint64, tx QueryerExecer) error {
	execer := r.getExecer(tx)

	query := sq.Delete("service_events").PlaceholderFormat(sq.Dollar)
	query = query.Where(sq.Eq{"id": eventIDs})

	s, args, err := query.ToSql()

	if err != nil {
		return err
	}

	_, err = execer.ExecContext(ctx, s, args...)

	if err == nil {
		metrics.SubEventsCountInPool(uint(len(eventIDs)))
	}

	return err
}
