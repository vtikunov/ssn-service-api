package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
)

func (o *serviceAPI) RemoveServiceV1(ctx context.Context, req *pb.RemoveServiceV1Request) (*pb.RemoveServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "RemoveServiceV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := o.srvService.Remove(ctx, req.ServiceId)
	if err != nil {
		logger.ErrorKV(ctx, "RemoveServiceV1 - failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveServiceV1Response{}, nil
}
