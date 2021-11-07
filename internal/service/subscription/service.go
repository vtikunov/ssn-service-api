package subscription

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

type serviceRepo interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	Add(ctx context.Context, service *subscription.Service) error
	Update(ctx context.Context, service *subscription.Service) error
	List(ctx context.Context) ([]*subscription.Service, error)
	Remove(ctx context.Context, serviceID uint64) (ok bool, err error)
}

type serviceService struct {
	repo serviceRepo
}

// NewServiceService создаёт инстанс сервиса ServiceService
func NewServiceService(repo serviceRepo) *serviceService {
	return &serviceService{
		repo: repo,
	}
}

// Describe - возвращает сервис по его ID.
func (s *serviceService) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	return s.repo.Describe(ctx, serviceID)
}

// Add - добавляет сервис.
func (s *serviceService) Add(ctx context.Context, service *subscription.Service) error {
	return s.repo.Add(ctx, service)
}

// Update - обновляет сервис.
func (s *serviceService) Update(ctx context.Context, service *subscription.Service) error {
	return s.repo.Update(ctx, service)
}

// List - возвращает постраничный список сервисов.
func (s *serviceService) List(ctx context.Context) ([]*subscription.Service, error) {
	return s.repo.List(ctx)
}

// Remove - удаляет сервис.
// Возвращает true если сервис существовал и успешно удален методом.
func (s serviceService) Remove(ctx context.Context, serviceID uint64) (ok bool, err error) {
	return s.repo.Remove(ctx, serviceID)
}
