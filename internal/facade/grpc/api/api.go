package api

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"

	pbf "github.com/ozonmp/ssn-service-api/pkg/ssn-service-facade"
)

type serviceService interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	List(ctx context.Context, offset uint64, limit uint64) ([]*subscription.Service, error)
}

type serviceAPI struct {
	pbf.UnimplementedSsnServiceFacadeServiceServer
	srvService serviceService
}

// NewServiceAPI returns api of ssn-service-api service
func NewServiceAPI(srv serviceService) pbf.SsnServiceFacadeServiceServer {
	return &serviceAPI{srvService: srv}
}

func convertServiceToPb(service *subscription.Service) *pbf.Service {
	return &pbf.Service{
		Id:          service.ID,
		Name:        service.Name,
		Description: service.Description,
	}
}
