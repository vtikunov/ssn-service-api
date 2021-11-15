package repo

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/retranslator/metrics"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
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
	query = query.Where(sq.Eq{"status": subscription.Deferred})
	query = query.Suffix("RETURNING id, service_id, type, status, payload, updated_at")

	s, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := execer.QueryContext(ctx, s, args...)
	defer func() {
		if rows == nil {
			return
		}

		if errCl := rows.Close(); errCl != nil {
			logger.ErrorKV(ctx, "eventRepo.Lock - failed close rows", "err", errCl)
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

	metrics.AddEventsCountInPool(uint(len(res)))

	return res, err
}
