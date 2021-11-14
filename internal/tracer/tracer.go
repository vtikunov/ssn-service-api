package tracer

import (
	"context"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

// NewTracer - returns new tracer.
func NewTracer(ctx context.Context, serviceName string, host string, port string) (io.Closer, error) {
	cfgTracer := &jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: host + port,
		},
	}
	tracer, closer, err := cfgTracer.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		logger.ErrorKV(ctx, "failed init jaeger", "err", err)

		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	logger.InfoKV(ctx, "traces started")

	return closer, nil
}
