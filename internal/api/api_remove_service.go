package api

import (
	"context"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *serviceAPI) RemoveServiceV1(
	ctx context.Context,
	req *pb.RemoveServiceV1Request,
) (*pb.RemoveServiceV1Response, error) {

	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("RemoveServiceV1 - invalid argument")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err := o.repo.Remove(ctx, req.ServiceId)
	if err != nil {
		log.Error().Err(err).Msg("RemoveServiceV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("RemoveServiceV1 - success")

	return &pb.RemoveServiceV1Response{
		IsFounded: ok,
	}, nil
}
