package service

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/bot/model/subscription"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
	pbf "github.com/ozonmp/ssn-service-api/pkg/ssn-service-facade"
)

type serviceReadClient struct {
	grpcClient pbf.SsnServiceFacadeServiceClient
}

// NewServiceReadClient - creating grpc client reader of services.
func NewServiceReadClient(grpcClient pbf.SsnServiceFacadeServiceClient) *serviceReadClient {
	return &serviceReadClient{
		grpcClient: grpcClient,
	}
}

func (s *serviceReadClient) Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error) {
	res, err := s.grpcClient.DescribeServiceV1(ctx, &pbf.DescribeServiceV1Request{ServiceId: serviceID})

	if err != nil {
		return nil, err
	}

	return convertPbReadService(res.GetService()), nil
}

func (s *serviceReadClient) List(ctx context.Context, offset uint64, limit uint64) (*subscription.ServiceListChunk, error) {
	res, err := s.grpcClient.ListServicesV1(ctx, &pbf.ListServicesV1Request{
		Offset: offset,
		Limit:  limit,
	})

	if err != nil {
		return nil, err
	}

	services := make([]*subscription.Service, len(res.GetServices()))

	for k, s := range res.GetServices() {
		services[k] = convertPbReadService(s)
	}

	var prevServices subscription.ServiceListChunkGetter
	var nextServices subscription.ServiceListChunkGetter

	if res.IsHasPrevPage {
		prevServices = s.makeServiceListChunkGetter(ctx, offset-limit, limit)
	}

	if res.IsHasNextPage {
		nextServices = s.makeServiceListChunkGetter(ctx, offset+limit, limit)
	}

	return &subscription.ServiceListChunk{
		Services:         services,
		PreviousServices: prevServices,
		NextServices:     nextServices,
	}, nil
}

func (s *serviceReadClient) makeServiceListChunkGetter(ctx context.Context, offset uint64, limit uint64) subscription.ServiceListChunkGetter {
	return func() (*subscription.ServiceListChunk, error) {
		return s.List(ctx, offset, limit)
	}
}

func convertPbReadService(service *pbf.Service) *subscription.Service {
	return &subscription.Service{
		ID:          service.GetId(),
		Name:        service.GetName(),
		Description: service.GetDescription(),
	}
}

type serviceWriteClient struct {
	grpcClient pb.SsnServiceApiServiceClient
}

// NewServiceWriteClient - creating grpc client writer of services.
func NewServiceWriteClient(grpcClient pb.SsnServiceApiServiceClient) *serviceWriteClient {
	return &serviceWriteClient{
		grpcClient: grpcClient,
	}
}

func (s *serviceWriteClient) Create(ctx context.Context, service *subscription.Service) (uint64, error) {
	res, err := s.grpcClient.CreateServiceV1(ctx, &pb.CreateServiceV1Request{
		Name:        service.Name,
		Description: service.Description,
	})

	if err != nil {
		return 0, err
	}

	return res.GetServiceId(), nil
}

func (s *serviceWriteClient) Update(ctx context.Context, serviceID uint64, service *subscription.Service) error {
	_, err := s.grpcClient.UpdateServiceV1(ctx, &pb.UpdateServiceV1Request{
		ServiceId:   serviceID,
		Name:        service.Name,
		Description: service.Description,
	})

	return err
}

func (s *serviceWriteClient) Remove(ctx context.Context, serviceID uint64) (bool, error) {
	_, err := s.grpcClient.RemoveServiceV1(ctx, &pb.RemoveServiceV1Request{
		ServiceId: serviceID,
	})

	return err == nil, err
}
