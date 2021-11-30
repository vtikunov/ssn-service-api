package subscription

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/facade/metrics"
	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/facade/repo"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"github.com/ozonmp/ssn-service-api/internal/tracer"
)

type serviceRepo interface {
	Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error)
	Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) (bool, error)
	Remove(ctx context.Context, serviceID, eventID uint64, tx repo.QueryerExecer) (bool, error)
	Describe(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (*subscription.Service, error)
	List(ctx context.Context, offset uint64, limit uint64, tx repo.QueryerExecer) ([]*subscription.Service, error)
}

type serviceCache interface {
	Set(ctx context.Context, service *subscription.Service) error
	Unset(ctx context.Context, serviceID uint64) error
	Get(ctx context.Context, serviceID uint64) (*subscription.Service, error)
}

type transactionalSession interface {
	Execute(ctx context.Context, fn func(ctx context.Context, tx repo.QueryerExecer) error) error
}

type serviceService struct {
	srvRepo  serviceRepo
	txs      transactionalSession
	srvCache serviceCache
}

// NewServiceService создаёт инстанс сервиса ServiceService
func NewServiceService(srvRepo serviceRepo, txs transactionalSession, srvCache serviceCache) *serviceService {
	return &serviceService{
		srvRepo:  srvRepo,
		txs:      txs,
		srvCache: srvCache,
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

		if err != nil {
			return err
		}

		isAdded = ok
		return nil
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
		if err != nil {
			return err
		}

		if s.isCacheAware() {
			if err = s.srvCache.Unset(ctx, service.ID); err != nil {
				return err
			}
		}

		isUpdated = ok

		return nil
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
		if err != nil {
			return err
		}

		if s.isCacheAware() {
			if err = s.srvCache.Unset(ctx, serviceID); err != nil {
				return err
			}
		}

		isRemoved = ok

		return nil
	})

	if err == nil && isRemoved {
		metrics.AddCudCountTotal(1, subscription.Removed)
	}

	return err
}

func (s *serviceService) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "service.Describe")
	defer sp.Finish()

	if !s.isCacheAware() {
		return s.srvRepo.Describe(ctx, serviceID, nil)
	}

	service, err := s.srvCache.Get(ctx, serviceID)

	if err != nil {
		logger.ErrorKV(ctx, "service.Describe: failed getting service from cache", "err", err, "ID", serviceID)
	}

	if service != nil && err == nil {
		return service, nil
	}

	service, err = s.srvRepo.Describe(ctx, serviceID, nil)

	if service != nil && err == nil {
		if cacheErr := s.srvCache.Set(ctx, service); cacheErr != nil {
			logger.ErrorKV(ctx, "service.Describe: failed setting service in cache", "err", err, "service", service)
		}
	}

	return service, err
}

func (s *serviceService) List(ctx context.Context, offset uint64, limit uint64) ([]*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "service.List")
	defer sp.Finish()

	return s.srvRepo.List(ctx, offset, limit, nil)
}

func (s *serviceService) isCacheAware() bool {
	return s.srvCache != nil
}
