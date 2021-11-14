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

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/config"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/repo"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/sender"

	retranslatorpkg "github.com/ozonmp/ssn-service-api/internal/retranslator/retranslator"
)

type retranslatorServer struct {
	db *sqlx.DB
}

// NewRetranslatorServer - создает сервис-окружение для запуска ретранслятора.
func NewRetranslatorServer(db *sqlx.DB) *retranslatorServer {
	return &retranslatorServer{db: db}
}

func (s *retranslatorServer) Start(ctx context.Context, cfg *config.Config) {
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

	retranslator := retranslatorpkg.NewRetranslator(
		&cfg.Retranslator,
		repo.NewEventRepo(s.db),
		sender.NewDummySender(),
	)

	retranslator.Start(ctx)

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

	retranslator.StopWait()

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

	logger.InfoKV(ctx, "retranslator shut down correctly")
}
