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

func (o *serviceAPI) CreateServiceV1(ctx context.Context, req *pb.CreateServiceV1Request) (*pb.CreateServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		logger.WarnKV(ctx, "CreateServiceV1 - invalid argument", "err", err)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	service := subscription.Service{Name: req.Name, Description: req.Description, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if err := o.srvService.Add(ctx, &service); err != nil {
		logger.ErrorKV(ctx, "CreateServiceV1 - failed", "err", err)

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateServiceV1Response{
		ServiceId: service.ID,
	}, nil
}
