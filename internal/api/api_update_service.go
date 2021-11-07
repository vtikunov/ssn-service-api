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

func (o *serviceAPI) UpdateServiceV1(
	ctx context.Context,
	req *pb.UpdateServiceV1Request,
) (*pb.UpdateServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("UpdateServiceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	service := subscription.Service{ID: req.ServiceId, Name: req.Name, Description: req.Description, UpdatedAt: time.Now()}

	if err := o.repo.Update(ctx, &service); err != nil {
		log.Error().Err(err).Msg("UpdateServiceV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("UpdateServiceV1 - success")

	return &pb.UpdateServiceV1Response{}, nil
}
