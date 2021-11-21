package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/ozonmp/ssn-service-api/internal/database"
	"github.com/ozonmp/ssn-service-api/internal/facade/config"
	"github.com/ozonmp/ssn-service-api/internal/facade/metrics"
	"github.com/ozonmp/ssn-service-api/internal/facade/server"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"github.com/ozonmp/ssn-service-api/internal/tracer"
)

func main() {
	ctx := context.Background()

	if err := config.ReadConfigYML("config.facade.yml"); err != nil {
		logger.FatalKV(ctx, "failed init configuration", "err", err)
	}

	cfg := config.GetConfigInstance()

	syncLogger := logger.InitLogger(ctx, cfg.Project.Debug, "service", cfg.Project.Name)
	defer syncLogger()

	migration := flag.Bool("migration", true, "Defines the migration start option")
	flag.Parse()

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

	if *migration {
		if err = goose.Up(db.DB, cfg.Database.Migrations); err != nil {
			logger.ErrorKV(ctx, "migration failed", "err", err)

			return
		}
	}

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

	server.NewConsumerServer(db).Start(ctx, &cfg)
}
