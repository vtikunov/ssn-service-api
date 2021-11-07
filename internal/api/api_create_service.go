package api

import (
	"context"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *serviceAPI) CreateServiceV1(
	ctx context.Context,
	req *pb.CreateServiceV1Request,
) (*pb.CreateServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("CreateServiceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	service := subscription.Service{Name: req.Name, Description: req.Description, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if err := o.repo.Add(ctx, &service); err != nil {
		log.Error().Err(err).Msg("CreateServiceV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("CreateServiceV1 - success")

	return &pb.CreateServiceV1Response{
		ServiceId: service.ID,
	}, nil
}
