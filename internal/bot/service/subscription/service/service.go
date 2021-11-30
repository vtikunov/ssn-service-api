package service

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/bot/metrics"
	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"
)

type serviceReader interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	List(ctx context.Context, offset uint64, limit uint64) (*subscription.ServiceListChunk, error)
}

type serviceWriter interface {
	Create(ctx context.Context, service *subscription.Service) (uint64, error)
	Update(ctx context.Context, serviceID uint64, service *subscription.Service) error
	Remove(ctx context.Context, serviceID uint64) (bool, error)
}

type serviceService struct {
	reader serviceReader
	writer serviceWriter
}

// NewServiceService - creating application service of services.
func NewServiceService(reader serviceReader, writer serviceWriter) *serviceService {
	return &serviceService{
		reader: reader,
		writer: writer,
	}
}

func (s *serviceService) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	return s.reader.Describe(ctx, serviceID)
}

func (s *serviceService) List(ctx context.Context, offset uint64, limit uint64) (*subscription.ServiceListChunk, error) {
	return s.reader.List(ctx, offset, limit)
}

func (s *serviceService) Create(ctx context.Context, service *subscription.Service) (uint64, error) {
	metrics.AddCudCountTotal(1, subscription.Created)
	return s.writer.Create(ctx, service)
}

func (s *serviceService) Update(ctx context.Context, serviceID uint64, service *subscription.Service) error {
	metrics.AddCudCountTotal(1, subscription.Updated)
	return s.writer.Update(ctx, serviceID, service)
}

func (s *serviceService) Remove(ctx context.Context, serviceID uint64) (bool, error) {
	metrics.AddCudCountTotal(1, subscription.Removed)
	return s.writer.Remove(ctx, serviceID)
}
