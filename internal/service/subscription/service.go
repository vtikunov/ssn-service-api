package subscription

import (
	"context"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/metrics"

	"github.com/ozonmp/ssn-service-api/internal/tracer"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/repo"
)

type serviceRepo interface {
	Describe(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (*subscription.Service, error)
	Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error
	Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error
	List(ctx context.Context, offset uint64, limit uint64, tx repo.QueryerExecer) ([]*subscription.Service, error)
	Remove(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) error
}

type eventRepo interface {
	Add(ctx context.Context, event *subscription.ServiceEvent, tx repo.QueryerExecer) error
}

type transactionalSession interface {
	Execute(ctx context.Context, fn func(ctx context.Context, tx repo.QueryerExecer) error) error
}

type serviceService struct {
	srvRepo   serviceRepo
	eventRepo eventRepo
	txs       transactionalSession
}

// NewServiceService создаёт инстанс сервиса ServiceService
func NewServiceService(srvRepo serviceRepo, eventRepo eventRepo, txs transactionalSession) *serviceService {
	return &serviceService{
		srvRepo:   srvRepo,
		eventRepo: eventRepo,
		txs:       txs,
	}
}

// Describe - возвращает сервис по его ID.
func (s *serviceService) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "service.Describe")
	defer sp.Finish()

	return s.srvRepo.Describe(ctx, serviceID, nil)
}

// Add - добавляет сервис.
// nolint:dupl
func (s *serviceService) Add(ctx context.Context, service *subscription.Service) error {
	sp := tracer.StartSpanFromContext(ctx, "service.Add")
	defer sp.Finish()

	var addErr error

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		addErr = s.srvRepo.Add(ctx, service, tx)

		if addErr != nil {
			return addErr
		}

		return s.eventRepo.Add(ctx, &subscription.ServiceEvent{
			ServiceID: service.ID,
			Type:      subscription.Created,
			Status:    subscription.Deferred,
			Service:   service,
			UpdatedAt: time.Now(),
		}, tx)
	})

	if addErr != nil {
		return addErr
	}

	if err == nil {
		metrics.AddCudCountTotal(1, subscription.Created)
	}

	return err
}

// Update - обновляет сервис.
// nolint:dupl
func (s *serviceService) Update(ctx context.Context, service *subscription.Service) error {
	sp := tracer.StartSpanFromContext(ctx, "service.Update")
	defer sp.Finish()

	var updErr error

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		updErr = s.srvRepo.Update(ctx, service, tx)

		if updErr != nil {
			return updErr
		}

		return s.eventRepo.Add(ctx, &subscription.ServiceEvent{
			ServiceID: service.ID,
			Type:      subscription.Updated,
			Status:    subscription.Deferred,
			Service:   service,
			UpdatedAt: time.Now(),
		}, tx)
	})

	if updErr != nil {
		return updErr
	}

	if err == nil {
		metrics.AddCudCountTotal(1, subscription.Updated)
	}

	return err
}

// UpdateName - обновляет наименование сервиса.
// nolint:dupl
func (s *serviceService) UpdateName(ctx context.Context, serviceID uint64, name string) error {
	sp := tracer.StartSpanFromContext(ctx, "service.UpdateName")
	defer sp.Finish()

	var updErr error

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		srv, err := s.srvRepo.Describe(ctx, serviceID, tx)

		if err != nil {
			return err
		}

		srv.Name = name

		updErr = s.srvRepo.Update(ctx, srv, tx)

		if updErr != nil {
			return updErr
		}

		return s.eventRepo.Add(ctx, &subscription.ServiceEvent{
			ServiceID: srv.ID,
			Type:      subscription.Updated,
			SubType:   subscription.NameSubtype,
			Status:    subscription.Deferred,
			Service:   srv,
			UpdatedAt: time.Now(),
		}, tx)
	})

	if updErr != nil {
		return updErr
	}

	if err == nil {
		metrics.AddCudCountTotal(1, subscription.Updated)
	}

	return err
}

// UpdateDescription - обновляет описание сервиса.
// nolint:dupl
func (s *serviceService) UpdateDescription(ctx context.Context, serviceID uint64, desc string) error {
	sp := tracer.StartSpanFromContext(ctx, "service.UpdateDescription")
	defer sp.Finish()

	var updErr error

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		srv, err := s.srvRepo.Describe(ctx, serviceID, tx)

		if err != nil {
			return err
		}

		srv.Description = desc

		updErr = s.srvRepo.Update(ctx, srv, tx)

		if updErr != nil {
			return updErr
		}

		return s.eventRepo.Add(ctx, &subscription.ServiceEvent{
			ServiceID: srv.ID,
			Type:      subscription.Updated,
			SubType:   subscription.DescriptionSubType,
			Status:    subscription.Deferred,
			Service:   srv,
			UpdatedAt: time.Now(),
		}, tx)
	})

	if updErr != nil {
		return updErr
	}

	if err == nil {
		metrics.AddCudCountTotal(1, subscription.Updated)
	}

	return err
}

// List - возвращает постраничный список сервисов.
func (s *serviceService) List(ctx context.Context, offset uint64, limit uint64) ([]*subscription.Service, error) {
	sp := tracer.StartSpanFromContext(ctx, "service.List")
	defer sp.Finish()

	return s.srvRepo.List(ctx, offset, limit, nil)
}

// Remove - удаляет сервис.
func (s serviceService) Remove(ctx context.Context, serviceID uint64) error {
	sp := tracer.StartSpanFromContext(ctx, "service.Remove")
	defer sp.Finish()

	var rmvErr error

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		rmvErr = s.srvRepo.Remove(ctx, serviceID, tx)

		if rmvErr != nil {
			return rmvErr
		}

		return s.eventRepo.Add(ctx, &subscription.ServiceEvent{
			ServiceID: serviceID,
			Type:      subscription.Removed,
			Status:    subscription.Deferred,
			UpdatedAt: time.Now(),
		}, tx)
	})

	if rmvErr != nil {
		return rmvErr
	}

	if err == nil {
		metrics.AddCudCountTotal(1, subscription.Removed)
	}

	return err
}
