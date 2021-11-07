package repo

import "github.com/jmoiron/sqlx"

type QueryerExecer interface {
	sqlx.Execer
	sqlx.Queryer
	sqlx.QueryerContext
	sqlx.ExecerContext
}
