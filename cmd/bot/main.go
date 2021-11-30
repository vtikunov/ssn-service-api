package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/ozonmp/ssn-service-api/internal/bot/config"
	"github.com/ozonmp/ssn-service-api/internal/bot/metrics"
	"github.com/ozonmp/ssn-service-api/internal/bot/server"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"github.com/ozonmp/ssn-service-api/internal/tracer"
)

func main() {
	ctx := context.Background()

	_ = godotenv.Load()

	token, found := os.LookupEnv("TELEGRAM_TOKEN")
	if !found {
		logger.FatalKV(ctx, "environment variable TELEGRAM_TOKEN not found in .env")
	}

	if err := config.ReadConfigYML("config.bot.yml"); err != nil {
		logger.FatalKV(ctx, "failed init configuration", "err", err)
	}

	cfg := config.GetConfigInstance()

	syncLogger := logger.InitLogger(ctx, cfg.Project.Debug, "service", cfg.Project.Name)
	defer syncLogger()

	logger.InfoKV(ctx, fmt.Sprintf("Starting service: %s", cfg.Project.Name),
		"version", cfg.Project.Version,
		"commitHash", cfg.Project.CommitHash,
		"debug", cfg.Project.Debug,
		"environment", cfg.Project.Environment,
	)

	metrics.InitMetrics(&cfg)

	tracing, err := tracer.NewTracer(ctx, cfg.Jaeger.Service, cfg.Jaeger.Host, cfg.Jaeger.Port)
	if err != nil {
		logger.ErrorKV(ctx, "failed init tracing", "err", err)

		return
	}
	defer func() {
		if err := tracing.Close(); err != nil {
			logger.ErrorKV(ctx, "failed close tracer", "err", err)
		}
	}()

	server.NewBotServer().Start(ctx, &cfg, token)
}
