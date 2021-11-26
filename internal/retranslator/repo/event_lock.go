package repo

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/metrics"
)

var (
	eventLocker = int32(binary.BigEndian.Uint32([]byte("service_events")))
	lockLocker  = int32(binary.BigEndian.Uint32([]byte("lock")))
	sqlLock     = fmt.Sprintf("select pg_try_advisory_xact_lock(%d, %d)", eventLocker, lockLocker)
)

func (r *eventRepo) Lock(ctx context.Context, n uint64, tx QueryerExecer) ([]subscription.ServiceEvent, error) {
	execer := r.getExecer(tx)

	_, err := execer.ExecContext(ctx, sqlLock)
	if err != nil {
		return nil, err
	}

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

	if err == nil {
		metrics.AddEventsCountInPool(uint(len(res)))
	}

	return res, err
}
