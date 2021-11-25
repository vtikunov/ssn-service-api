package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/jmoiron/sqlx"

	"github.com/ozonmp/ssn-service-api/internal/facade/config"
	"github.com/ozonmp/ssn-service-api/internal/facade/database"
	"github.com/ozonmp/ssn-service-api/internal/facade/kafka"
	"github.com/ozonmp/ssn-service-api/internal/facade/service/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	servicerepo "github.com/ozonmp/ssn-service-api/internal/facade/repo/subscription/service"
)

type consumerServer struct {
	db *sqlx.DB
}

// NewConsumerServer - создает сервис-окружение для запуска консьюмеров событий.
func NewConsumerServer(db *sqlx.DB) *consumerServer {
	return &consumerServer{
		db: db,
	}
}

func (s *consumerServer) Start(ctx context.Context, cfg *config.Config) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	metricsAddr := fmt.Sprintf("%s:%v", cfg.Metrics.Host, cfg.Metrics.Port)
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

	r := servicerepo.NewServiceRepo(s.db)
	txs := database.NewTransactionalSession(s.db)
	srv := subscription.NewServiceService(r, txs)

	for cn := uint8(0); cn < cfg.Kafka.PartitionFactor; cn++ {
		err := kafka.StartConsuming(ctx, cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Group, kafka.GetServiceEventConsume(srv))
		if err != nil {
			logger.FatalKV(ctx, "failed start consuming", "err", err)
		}
	}

	isReady.Store(true)

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

	logger.InfoKV(ctx, "facade consumers shut down correctly")
}
