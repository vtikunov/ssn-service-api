package api

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalServiceNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "ssn_service_api_service_not_found_total",
		Help: "Total number of services that were not found",
	})
)

type serviceService interface {
	Describe(ctx context.Context, serviceID uint64) (*subscription.Service, error)
	Add(ctx context.Context, service *subscription.Service) error
	Update(ctx context.Context, service *subscription.Service) error
	List(ctx context.Context) ([]*subscription.Service, error)
	Remove(ctx context.Context, serviceID uint64) (ok bool, err error)
}

type serviceAPI struct {
	pb.UnimplementedSsnServiceApiServiceServer
	srvService serviceService
}

// NewServiceAPI returns api of ssn-service-api service
func NewServiceAPI(srv serviceService) pb.SsnServiceApiServiceServer {
	return &serviceAPI{srvService: srv}
}

func convertServiceToPb(service *subscription.Service) *pb.Service {
	return &pb.Service{
		Id:          service.ID,
		Name:        service.Name,
		Description: service.Description,
		CreatedAt:   timestamppb.New(service.CreatedAt),
		UpdatedAt:   timestamppb.New(service.UpdatedAt),
	}
}
