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

	ssn_service_facade "github.com/ozonmp/ssn-service-api/pkg/ssn-service-facade"

	ssn_service_api "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"

	"github.com/ozonmp/ssn-service-api/internal/bot/router"
	"github.com/ozonmp/ssn-service-api/internal/bot/service/subscription/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/ozonmp/ssn-service-api/internal/bot/config"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

type botServer struct {
}

// NewBotServer - создает сервис-окружение для запуска бота.
func NewBotServer() *botServer {
	return &botServer{}
}

func (s *botServer) Start(ctx context.Context, cfg *config.Config, token string) {
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

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.FatalKV(ctx, "failed connect to telegram api", "err", err)
	}

	logger.InfoKV(ctx, "bot authorized on account", "account", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Timeout: int(cfg.Bot.Timeout),
	}

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.FatalKV(ctx, "failed getting telegram api updates channel", "err", err)
	}

	wServiceServiceConn, err := grpc.DialContext(
		ctx,
		cfg.Bot.WriteServiceServiceAddr,
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(),
		),
	)
	if err != nil {
		logger.FatalKV(ctx, "failed connect to write service grpc server", "err", err)
	}

	rServiceServiceConn, err := grpc.DialContext(
		ctx,
		cfg.Bot.ReadServiceServiceAddr,
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(),
		),
	)
	if err != nil {
		logger.FatalKV(ctx, "failed connect to read service grpc server", "err", err)
	}

	srvService := service.NewServiceService(
		service.NewServiceReadClient(
			ssn_service_facade.NewSsnServiceFacadeServiceClient(rServiceServiceConn),
		),
		service.NewServiceWriteClient(
			ssn_service_api.NewSsnServiceApiServiceClient(wServiceServiceConn),
		),
	)

	r := router.NewRouter(bot, srvService, &cfg.Bot)

	go func() {
		for update := range updates {
			r.HandleUpdate(ctx, update)
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

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

	logger.InfoKV(ctx, "bot shut down correctly")
}
