package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/rs/zerolog/log"
)

func (r *eventRepo) Lock(ctx context.Context, n uint64, tx QueryerExecer) ([]subscription.ServiceEvent, error) {
	execer := r.getExecer(tx)

	subQ := sq.Select("s.id").From("service_events s").PlaceholderFormat(sq.Dollar)
	subQ = subQ.LeftJoin("service_events s1 on (s1.service_id = s.service_id AND s1.status = ?)", subscription.Processed)
	subQ = subQ.Where(sq.And{sq.Expr("s1.service_id IS NULL"), sq.Eq{"s.status": subscription.Deferred}})
	subQ = subQ.OrderBy("s.id").Limit(n)

	query := sq.Update("service_events").PlaceholderFormat(sq.Dollar)
	query = query.Set("status", subscription.Processed)
	query = query.Set("updated_at", "NOW()")
	query = query.Where(subQ.Prefix("id IN (").Suffix(")"))
	query = query.Suffix("RETURNING s.id, s.service_id, s.type, s.status, s.payload, s.updated_at")

	s, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := execer.QueryContext(ctx, s, args...)
	defer func() {
		if errCl := rows.Close(); errCl != nil {
			log.Error().Err(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	srvEvents := make([]*serviceEvent, 0)
	err = sqlx.StructScan(rows, &srvEvents)

	if err != nil {
		return nil, err
	}

	res := make([]subscription.ServiceEvent, 0, len(srvEvents))
	for _, event := range srvEvents {
		res = append(res, *event.convertToServiceEvent())
	}

	return res, err
}
