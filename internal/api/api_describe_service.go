package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
)

func (o *serviceAPI) DescribeServiceV1(ctx context.Context, req *pb.DescribeServiceV1Request) (*pb.DescribeServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "DescribeServiceV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	service, err := o.srvService.Describe(ctx, req.ServiceId)
	if err != nil {
		logger.ErrorKV(ctx, "DescribeServiceV1 - failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if service == nil {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	return &pb.DescribeServiceV1Response{
		Service: convertServiceToPb(service),
	}, nil
}
