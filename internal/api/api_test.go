package api

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/ozonmp/ssn-service-api/internal/service/subscription"

	apimocks "github.com/ozonmp/ssn-service-api/internal/mocks/api"
	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
)

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	ctrl := gomock.NewController(t)
	repo := apimocks.NewMockServiceRepo(ctrl)
	eventRepo := apimocks.NewMockEventRepo(ctrl)
	tsx := apimocks.NewMockTransactionalSession(ctrl)

	pb.RegisterSsnServiceApiServiceServer(server, NewServiceAPI(subscription.NewServiceService(repo, eventRepo, tsx)))

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func prepareClient(ctx context.Context, t *testing.T) (client pb.SsnServiceApiServiceClient, closeClient func()) {

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(t)))
	if err != nil {
		log.Fatal(err)
	}
	closeCl := func() {
		err := conn.Close()
		if err != nil {
			log.Panicln(err)
		}
	}

	return pb.NewSsnServiceApiServiceClient(conn), closeCl
}

//nolint:dupl
func TestServiceAPI_CreateServiceV1Request_NameValidation(t *testing.T) {
	ctx := context.Background()
	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.CreateServiceV1Request{
		{},
		{Name: ""},
	}

	for _, request := range requests {
		response, err := client.CreateServiceV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid CreateServiceV1Request.Name: value length must be between 1 and 100 runes, inclusive", er.Message())
	}
}

//nolint:dupl
func TestServiceAPI_DescribeServiceV1Request_ServiceIDValidation(t *testing.T) {
	ctx := context.Background()
	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.DescribeServiceV1Request{
		{},
		{ServiceId: 0},
	}

	for _, request := range requests {
		response, err := client.DescribeServiceV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid DescribeServiceV1Request.ServiceId: value must be greater than 0", er.Message())
	}
}

//nolint:dupl
func TestServiceAPI_ListServiceV1Request_LimitValidation(t *testing.T) {
	ctx := context.Background()
	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.ListServicesV1Request{
		{Offset: 10},
		{Offset: 10, Limit: 0},
		{Offset: 10, Limit: 501},
	}

	for _, request := range requests {
		response, err := client.ListServicesV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid ListServicesV1Request.Limit: value must be inside range (0, 500]", er.Message())
	}
}

//nolint:dupl
func TestServiceAPI_RemoveServiceV1Request_ServiceIDValidation(t *testing.T) {
	ctx := context.Background()
	client, closeCl := prepareClient(ctx, t)
	defer closeCl()

	requests := []*pb.RemoveServiceV1Request{
		{},
		{ServiceId: 0},
	}

	for _, request := range requests {
		response, err := client.RemoveServiceV1(ctx, request)

		assert.Nil(t, response)
		assert.NotNil(t, err)

		er, _ := status.FromError(err)

		assert.Equal(t, codes.InvalidArgument, er.Code())
		assert.Equal(t, "invalid RemoveServiceV1Request.ServiceId: value must be greater than 0", er.Message())
	}
}
