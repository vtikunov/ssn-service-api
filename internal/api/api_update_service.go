package api

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
)

func (o *serviceAPI) UpdateServiceV1(ctx context.Context, req *pb.UpdateServiceV1Request) (*pb.UpdateServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "UpdateServiceV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	service := subscription.Service{ID: req.ServiceId, Name: req.Name, Description: req.Description, UpdatedAt: time.Now()}

	if err := o.srvService.Update(ctx, &service); err != nil {
		logger.ErrorKV(ctx, "UpdateServiceV1 - failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateServiceV1Response{}, nil
}
