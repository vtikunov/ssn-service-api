package repo

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

var errNotImplementedMethod = errors.New("method is not implemented")

// ServiceRepo - DAO Service
type ServiceRepo interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	Add(ctx context.Context, service *subscription.Service) error
	List(ctx context.Context) ([]*subscription.Service, error)
	Remove(ctx context.Context, serviceID uint64) (ok bool, err error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

func (r repo) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	return nil, errNotImplementedMethod
}

func (r repo) Add(ctx context.Context, service *subscription.Service) error {
	return errNotImplementedMethod
}

func (r repo) List(ctx context.Context) ([]*subscription.Service, error) {
	return nil, errNotImplementedMethod
}

func (r repo) Remove(ctx context.Context, serviceID uint64) (ok bool, err error) {
	return false, errNotImplementedMethod
}

// NewRepo returns ServiceRepo interface
func NewRepo(db *sqlx.DB, batchSize uint) *repo {
	return &repo{db: db, batchSize: batchSize}
}
