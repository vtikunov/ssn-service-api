package grpc_logs

import (
	"context"
	"strings"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MetadataChangingLogsLevelUnaryServerInterceptor - добавляет возможность смены уровня логов через мета - Logs-Level.
func MetadataChangingLogsLevelUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			levels := md.Get("Log-Level")
			if len(levels) > 0 {
				logger.InfoKV(ctx, "got log level", "levels", levels)

				var isParsed bool

				if len(levels) == 1 {
					if parsedLevel, ok := parseLevel(levels[0]); ok {
						isParsed = ok
						newLogger := logger.CloneWithLevel(ctx, parsedLevel)
						ctx = logger.AttachLogger(ctx, newLogger)
						logger.InfoKV(ctx, "log level was changed", "level", levels[0])
					}
				}

				if !isParsed {
					logger.InfoKV(ctx, "log level was not parsed & changed")
				}
			}
		}

		return handler(ctx, req)
	}
}

func parseLevel(str string) (int8, bool) {
	switch strings.ToLower(str) {
	case "debug":
		return int8(logger.DebugLevel), true
	case "info":
		return int8(logger.InfoLevel), true
	case "warn":
		return int8(logger.WarnLevel), true
	case "error":
		return int8(logger.ErrorLevel), true
	default:
		return int8(logger.DebugLevel), false
	}
}
