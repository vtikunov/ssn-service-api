package subscription

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/repo"
)

type serviceRepo interface {
	Describe(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (*subscription.Service, error)
	Add(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error
	Update(ctx context.Context, service *subscription.Service, tx repo.QueryerExecer) error
	List(ctx context.Context, tx repo.QueryerExecer) ([]*subscription.Service, error)
	Remove(ctx context.Context, serviceID uint64, tx repo.QueryerExecer) (ok bool, err error)
}

type eventRepo interface {
	Add(ctx context.Context, event []subscription.ServiceEvent, tx repo.QueryerExecer) error
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
	return s.srvRepo.Describe(ctx, serviceID, nil)
}

// Add - добавляет сервис.
func (s *serviceService) Add(ctx context.Context, service *subscription.Service) error {
	var addErr error

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		addErr = s.srvRepo.Add(ctx, service, tx)

		if addErr != nil {
			return addErr
		}

		return s.eventRepo.Add(ctx, []subscription.ServiceEvent{
			{
				ServiceID: service.ID,
				Type:      subscription.Created,
				Service:   service,
			},
		}, tx)
	})

	if addErr != nil {
		return addErr
	}

	return err
}

// Update - обновляет сервис.
func (s *serviceService) Update(ctx context.Context, service *subscription.Service) error {
	var updErr error

	err := s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		updErr = s.srvRepo.Update(ctx, service, tx)

		if updErr != nil {
			return updErr
		}

		return s.eventRepo.Add(ctx, []subscription.ServiceEvent{
			{
				ServiceID: service.ID,
				Type:      subscription.Updated,
				Service:   service,
			},
		}, tx)
	})

	if updErr != nil {
		return updErr
	}

	return err
}

// List - возвращает постраничный список сервисов.
func (s *serviceService) List(ctx context.Context) ([]*subscription.Service, error) {
	return s.srvRepo.List(ctx, nil)
}

// Remove - удаляет сервис.
// Возвращает true если сервис существовал и успешно удален методом.
func (s serviceService) Remove(ctx context.Context, serviceID uint64) (ok bool, err error) {
	var rmvErr error
	var rmvOk bool

	err = s.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
		rmvOk, rmvErr = s.srvRepo.Remove(ctx, serviceID, tx)

		if rmvErr != nil {
			return rmvErr
		}

		if rmvOk {
			return s.eventRepo.Add(ctx, []subscription.ServiceEvent{
				{
					ServiceID: serviceID,
					Type:      subscription.Removed,
				},
			}, tx)
		}

		return nil
	})

	if rmvErr != nil || err == nil {
		return rmvOk, rmvErr
	}

	return false, err
}
