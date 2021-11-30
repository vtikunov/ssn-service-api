package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	pbf "github.com/ozonmp/ssn-service-api/pkg/ssn-service-facade"
)

func (o *serviceAPI) ListServicesV1(ctx context.Context, req *pbf.ListServicesV1Request) (*pbf.ListServicesV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "ListServicesV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var isHasPrevPage, isHasNextPage bool

	if req.Offset > 0 {
		isHasPrevPage = true
	}

	services, err := o.srvService.List(ctx, req.Offset, req.Limit+1)
	if err != nil {
		logger.ErrorKV(ctx, "ListServicesV1 - failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if len(services) > int(req.Limit) {
		isHasNextPage = true
		services = services[:len(services)-1]
	}

	servicesPb := make([]*pbf.Service, 0, len(services))

	for _, service := range services {
		servicesPb = append(servicesPb, convertServiceToPb(service))
	}

	return &pbf.ListServicesV1Response{
		Services:      servicesPb,
		IsHasPrevPage: isHasPrevPage,
		IsHasNextPage: isHasNextPage,
	}, nil
}
