package grpc_logs

import (
	"context"

	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

// GetIsEnableDescribeRequestAndResponseDecider - создает decider для интерцептора zap логгера.
func GetIsEnableDescribeRequestAndResponseDecider() grpc_logging.ServerPayloadLoggingDecider {
	return func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return logger.FromContext(ctx).Desugar().Core().Enabled(logger.DebugLevel)
	}
}
