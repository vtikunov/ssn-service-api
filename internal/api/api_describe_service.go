package api

import (
	"context"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *serviceAPI) DescribeServiceV1(
	ctx context.Context,
	req *pb.DescribeServiceV1Request,
) (*pb.DescribeServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("DescribeServiceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	service, err := o.srvService.Describe(ctx, req.ServiceId)
	if err != nil {
		log.Error().Err(err).Msg("DescribeServiceV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if service == nil {
		log.Debug().Uint64("serviceId", req.ServiceId).Msg("service not found")
		totalServiceNotFound.Inc()

		return nil, status.Error(codes.NotFound, "service not found")
	}

	log.Debug().Msg("DescribeServiceV1 - success")

	return &pb.DescribeServiceV1Response{
		Service: convertServiceToPb(service),
	}, nil
}
