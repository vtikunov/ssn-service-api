package main

import (
	"context"
	"fmt"

	"github.com/ozonmp/ssn-service-api/internal/retranslator/sender"

	"github.com/ozonmp/ssn-service-api/internal/retranslator/metrics"

	"github.com/ozonmp/ssn-service-api/internal/retranslator/server"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/ozonmp/ssn-service-api/internal/database"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/config"
	"github.com/ozonmp/ssn-service-api/internal/tracer"
)

func main() {
	ctx := context.Background()

	if err := config.ReadConfigYML("config.retranslator.yml"); err != nil {
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

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	db, err := database.NewPostgres(ctx, dsn, cfg.Database.Driver)
	if err != nil {
		logger.FatalKV(ctx, "failed init postgres", "err", err)
	}
	defer func() {
		if errCl := db.Close(); errCl != nil {
			logger.ErrorKV(ctx, "failed close DB connection", "err", errCl)
		}
	}()

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

	snd := sender.NewKafkaSender(ctx, cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.PartitionFactor, cfg.Kafka.SendRetryMax)

	server.NewRetranslatorServer(db, snd).Start(ctx, &cfg)
}
