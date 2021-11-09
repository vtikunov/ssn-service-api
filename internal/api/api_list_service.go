package api

import (
	"context"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *serviceAPI) ListServicesV1(
	ctx context.Context,
	req *pb.ListServicesV1Request,
) (*pb.ListServicesV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("ListServicesV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	services, err := o.srvService.List(ctx, req.Offset, req.Limit)
	if err != nil {
		log.Error().Err(err).Msg("ListServicesV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("ListServicesV1 - success")

	servicesPb := make([]*pb.Service, 0, len(services))

	for _, service := range services {
		servicesPb = append(servicesPb, convertServiceToPb(service))
	}

	return &pb.ListServicesV1Response{
		Services: servicesPb,
	}, nil
}
