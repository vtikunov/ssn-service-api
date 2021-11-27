package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/jmoiron/sqlx"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/ozonmp/ssn-service-api/internal/facade/database"
	"github.com/ozonmp/ssn-service-api/internal/facade/grpc/api"
	"github.com/ozonmp/ssn-service-api/internal/facade/grpc/config"
	"github.com/ozonmp/ssn-service-api/internal/facade/service/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/grpc/interceptor/grpc_logs"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	servicerepo "github.com/ozonmp/ssn-service-api/internal/facade/repo/subscription/service"
	pbf "github.com/ozonmp/ssn-service-api/pkg/ssn-service-facade"
)

type grpcServer struct {
	db *sqlx.DB
}

// NewGrpcServer - создает сервис grpc.
func NewGrpcServer(db *sqlx.DB) *grpcServer {
	return &grpcServer{
		db: db,
	}
}

func (s *grpcServer) Start(ctx context.Context, cfg *config.Config) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)
	grpcAddr := fmt.Sprintf("%s:%v", cfg.Grpc.Host, cfg.Grpc.Port)
	metricsServer := createMetricsServer(cfg)

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("metrics server is running on %s", metricsAddr))
		if err := metricsServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "failed running metrics server", "err", err)
			cancel()
		}
	}()

	isReady := &atomic.Value{}
	isReady.Store(false)

	statusServer := createStatusServer(ctx, cfg, isReady)

	go func() {
		statusAdrr := fmt.Sprintf("%s:%v", cfg.Status.Host, cfg.Status.Port)
		logger.InfoKV(ctx, fmt.Sprintf("status server is running on %s", statusAdrr))
		if err := statusServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorKV(ctx, "failed running status server", "err", err)
		}
	}()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.ErrorKV(ctx, "failed to listen: %w", "err", err)
		cancel()
		return
	}
	defer func() {
		if err := l.Close(); err != nil {
			logger.DebugKV(ctx, "failed close listen", "err", err)
		}
	}()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Duration(cfg.Grpc.MaxConnectionIdle) * time.Minute,
			Timeout:           time.Duration(cfg.Grpc.Timeout) * time.Second,
			MaxConnectionAge:  time.Duration(cfg.Grpc.MaxConnectionAge) * time.Minute,
			Time:              time.Duration(cfg.Grpc.Timeout) * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
			grpc_logs.MetadataChangingLogsLevelUnaryServerInterceptor(),
			grpc_zap.PayloadUnaryServerInterceptor(logger.FromContext(ctx).Desugar(), grpc_logs.GetIsEnableDescribeRequestAndResponseDecider()),
		)),
	)

	r := servicerepo.NewServiceRepo(s.db)
	txs := database.NewTransactionalSession(s.db)
	srv := subscription.NewServiceService(r, txs)

	pbf.RegisterSsnServiceFacadeServiceServer(grpcServer, api.NewServiceAPI(srv))

	go func() {
		logger.InfoKV(ctx, fmt.Sprintf("GRPC Server is listening on: %s", grpcAddr))
		if err := grpcServer.Serve(l); err != nil {
			logger.FatalKV(ctx, "failed running gRPC server", "err", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		isReady.Store(true)
		logger.InfoKV(ctx, "the service is ready to accept requests")
	}()

	if cfg.Project.Debug {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		logger.InfoKV(ctx, fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		logger.InfoKV(ctx, fmt.Sprintf("ctx.Done: %v", done))
	}

	isReady.Store(false)

	if err := statusServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "statusServer.Shutdown", "err", err)
	} else {
		logger.InfoKV(ctx, "statusServer shut down correctly")
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.ErrorKV(ctx, "metricsServer.Shutdown", "err", err)
	} else {
		logger.InfoKV(ctx, "metricsServer shut down correctly")
	}

	logger.InfoKV(ctx, "facade gRPC server shut down correctly")
}
