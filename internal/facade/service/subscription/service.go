package subscription

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/facade/metrics"
	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/facade/repo"
	"github.com/ozonmp/ssn-service-api/internal/tracer"
)

type serviceRepo interface {
	Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error)
	Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error)
	Remove(ctx context.Context, serviceID, eventID uint64, tx repo.QueryerExecer) (bool, error)
	Describe(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (*subscription.Service, error)
	List(ctx context.Context, offset uint64, limit uint64, tx repo.QueryerExecer) ([]*subscription.Service, error)
}

type transactionalSession interface {
	Execute(ctx context.Context, fn func(ctx context.Context, tx repo.QueryerExecer) error) error
}

type serviceService struct {
	srvRepo serviceRepo
	txs     transactionalSession
}

// NewServiceService создаёт инстанс сервиса ServiceService
func NewServiceService(srvRepo serviceRepo, txs transactionalSession) *serviceService {
	return &serviceService{
		srvRepo: srvRepo,
		txs:     txs,
	}
}

// Add - добавляет сервис.
// nolint:dupl
func (s *serviceService) Add(ctx context.Context, service *subscription.Service) error {
	sp := tracer.StartSpanFromContext(ctx, "service.Add")
	defer sp.Finish()

	var isAdded bool

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		ok, err := s.srvRepo.Add(ctx, service, tx)
		isAdded = ok

		return err
	})

	if err == nil && isAdded {
		metrics.AddCudCountTotal(1, subscription.Created)
	}

	return err
}

// Update - обновляет сервис.
// nolint:dupl
func (s *serviceService) Update(ctx context.Context, service *subscription.Service) error {
	sp := tracer.StartSpanFromContext(ctx, "service.Update")
	defer sp.Finish()

	var isUpdated bool

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		ok, err := s.srvRepo.Update(ctx, service, tx)
		isUpdated = ok

		return err
	})

	if err == nil && isUpdated {
		metrics.AddCudCountTotal(1, subscription.Updated)
	}

	return err
}

// Remove - удаляет сервис.
func (s *serviceService) Remove(ctx context.Context, serviceID uint64, eventID uint64) error {
	sp := tracer.StartSpanFromContext(ctx, "service.Remove")
	defer sp.Finish()

	var isRemoved bool

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		ok, err := s.srvRepo.Remove(ctx, serviceID, eventID, tx)
		isRemoved = ok

		return err
	})

	if err == nil && isRemoved {
		metrics.AddCudCountTotal(1, subscription.Removed)
	}

	return err
}

func (s *serviceService) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "service.Describe")
	defer sp.Finish()

	return s.srvRepo.Describe(ctx, serviceID, nil)
}

func (s *serviceService) List(ctx context.Context, offset uint64, limit uint64) ([]*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "service.List")
	defer sp.Finish()

	return s.srvRepo.List(ctx, offset, limit, nil)
}
