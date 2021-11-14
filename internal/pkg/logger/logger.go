package logger

import (
	"context"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var attachedLoggerKey = &ctxKey{}

var globalLogger *zap.SugaredLogger

const (
	// DebugLevel - DEBUG.
	DebugLevel = zapcore.DebugLevel
	// InfoLevel - INFO.
	InfoLevel = zapcore.InfoLevel
	// WarnLevel - WARNINGS.
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel - ERRORS.
	ErrorLevel = zapcore.ErrorLevel
)

// InitLogger - инициализирует логгер.
func InitLogger(ctx context.Context, debug bool, kvs ...interface{}) (syncFn func()) {
	loggingLevel := zap.InfoLevel
	if debug {
		loggingLevel = zap.DebugLevel
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zap.NewAtomicLevelAt(loggingLevel),
	)

	notSugaredLogger := zap.New(consoleCore)

	sugaredLogger := notSugaredLogger.Sugar()
	SetLogger(sugaredLogger.With(kvs...))

	return func() {
		if errSync := notSugaredLogger.Sync(); errSync != nil {
			ErrorKV(ctx, "failed to sync logger", "err", errSync)
		}
	}
}

// FromContext - получает текущий логгер.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	var result *zap.SugaredLogger
	if attachedLogger, ok := ctx.Value(attachedLoggerKey).(*zap.SugaredLogger); ok {
		result = attachedLogger
	} else {
		result = globalLogger
	}

	jaegerSpan := opentracing.SpanFromContext(ctx)
	if jaegerSpan != nil {
		if spanCtx, ok := opentracing.SpanFromContext(ctx).Context().(jaeger.SpanContext); ok {
			result = result.With("trace-id", spanCtx.TraceID())
		}
	}

	return result
}

// ErrorKV - логирует с уровнем Error.
func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Errorw(message, kvs...)
}

// WarnKV - логирует с уровнем Warning.
func WarnKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Warnw(message, kvs...)
}

// InfoKV - логирует с уровнем Info.
func InfoKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Infow(message, kvs...)
}

// DebugKV - логирует с уровнем Debug.
func DebugKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Debugw(message, kvs...)
}

// FatalKV - логирует с уровнем Fatal и завершает работу.
func FatalKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Fatalw(message, kvs...)
}

// AttachLogger - передача логгера в контекст.
func AttachLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, attachedLoggerKey, logger)
}

// CloneWithLevel - клонирует логгер с необходимым уровнем логирования.
func CloneWithLevel(ctx context.Context, newLevel int8) *zap.SugaredLogger {
	return FromContext(ctx).
		Desugar().
		WithOptions(WithLevel(zapcore.Level(newLevel))).
		Sugar()
}

// SetLogger - устанавливает глобальный логгер.
func SetLogger(newLogger *zap.SugaredLogger) {
	globalLogger = newLogger
}

func init() {
	notSugaredLogger, err := zap.NewProduction()
	if err != nil {
		log.Panic(err)
	}

	globalLogger = notSugaredLogger.Sugar()
}
