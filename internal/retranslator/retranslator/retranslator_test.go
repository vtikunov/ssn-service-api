package retranslator_test

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"go.uber.org/zap"

	"github.com/ozonmp/ssn-service-api/internal/retranslator/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"

	retranslatormocks "github.com/ozonmp/ssn-service-api/internal/mocks/retranslator"
	repopkg "github.com/ozonmp/ssn-service-api/internal/retranslator/repo"
	retranslatorpkg "github.com/ozonmp/ssn-service-api/internal/retranslator/retranslator"
)

type initData struct {
	eventChannelSize uint64

	maxConsumers      uint64
	consumerTimeout   time.Duration
	consumerBatchTime time.Duration
	consumerBatchSize uint64

	maxProducers       uint64
	producerTimeout    time.Duration
	producerMaxWorkers uint64

	isEmitSendError bool
}

func initLogger() {
	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		zap.NewAtomicLevelAt(zap.PanicLevel),
	)
	notSugaredLogger := zap.New(consoleCore)
	sugaredLogger := notSugaredLogger.Sugar()
	logger.SetLogger(sugaredLogger)
}

func SuiteAllEventsCompleteWhenStoppingByFunc(t *testing.T, d initData) {

	ctrl := gomock.NewController(t)
	ctx := context.Background()

	initLogger()

	repo := retranslatormocks.NewMockEventRepo(ctrl)
	sender := retranslatormocks.NewMockEventSender(ctrl)

	var lockCount int64
	//nolint:dupl
	repo.EXPECT().Lock(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64, tx repopkg.QueryerExecer) ([]subscription.ServiceEvent, error) {
			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i, ServiceID: start + i, Service: &subscription.Service{ID: start + i}}
			}

			atomic.AddInt64(&lockCount, int64(n))

			return result, nil
		},
	).AnyTimes()

	var unlockCount int64
	repo.EXPECT().Unlock(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, eventIDs []uint64, tx repopkg.QueryerExecer) error {
			atomic.AddInt64(&unlockCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var removeCount int64
	repo.EXPECT().Remove(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, eventIDs []uint64, tx repopkg.QueryerExecer) error {
			atomic.AddInt64(&removeCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var sendCount int64
	var sendErrorCount int64
	sender.EXPECT().Send(gomock.Any()).DoAndReturn(
		func(serviceEvent *subscription.ServiceEvent) error {
			if d.isEmitSendError {
				if serviceEvent.ID%2 == 0 {
					atomic.AddInt64(&sendErrorCount, 1)

					return errors.New("i can't send")
				}
			}
			atomic.AddInt64(&sendCount, 1)

			return nil
		},
	).AnyTimes()

	retranslator := retranslatorpkg.NewRetranslator(
		ctx,
		&config.Retranslator{
			EventChannelSize: d.eventChannelSize,

			MaxConsumers:      d.maxConsumers,
			ConsumerTimeout:   d.consumerTimeout,
			ConsumerBatchTime: d.consumerBatchTime,
			ConsumerBatchSize: d.consumerBatchSize,

			MaxProducers:       d.maxProducers,
			ProducerTimeout:    d.producerTimeout,
			ProducerMaxWorkers: d.producerMaxWorkers,
		},
		repo,
		sender,
	)

	retranslator.Start(ctx)

	time.Sleep(time.Millisecond * 500)

	retranslator.StopWait()

	if !d.isEmitSendError {
		assert.LessOrEqual(t, sendCount, lockCount)
		assert.Equal(t, sendErrorCount, unlockCount)

		return
	}

	assert.Equal(t, sendCount, removeCount)
	assert.Equal(t, lockCount, unlockCount+removeCount)
}

func TestAllEventsCompleteWhenStoppingByFuncOne(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			eventChannelSize: 1,

			maxConsumers:      10,
			consumerTimeout:   time.Millisecond,
			consumerBatchTime: time.Microsecond,
			consumerBatchSize: 10,

			maxProducers:       10,
			producerTimeout:    time.Millisecond,
			producerMaxWorkers: 2,

			isEmitSendError: false,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFuncTwo(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			eventChannelSize: 10,

			maxConsumers:      20,
			consumerTimeout:   time.Millisecond,
			consumerBatchTime: time.Microsecond,
			consumerBatchSize: 10,

			maxProducers:       20,
			producerTimeout:    time.Millisecond,
			producerMaxWorkers: 2,

			isEmitSendError: true,
		},
	)
}

func SuiteAllEventsCompleteWhenStoppingByContext(t *testing.T, d initData) {

	ctrl := gomock.NewController(t)
	ctx, cancelCtx := context.WithCancel(context.Background())

	initLogger()

	repo := retranslatormocks.NewMockEventRepo(ctrl)
	sender := retranslatormocks.NewMockEventSender(ctrl)

	var lockCount int64
	//nolint:dupl
	repo.EXPECT().Lock(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64, tx repopkg.QueryerExecer) ([]subscription.ServiceEvent, error) {
			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i, ServiceID: start + i, Service: &subscription.Service{ID: start + i}}
			}

			atomic.AddInt64(&lockCount, int64(n))

			return result, nil
		},
	).AnyTimes()

	var unlockCount int64
	repo.EXPECT().Unlock(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, eventIDs []uint64, tx repopkg.QueryerExecer) error {
			atomic.AddInt64(&unlockCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var removeCount int64
	repo.EXPECT().Remove(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, eventIDs []uint64, tx repopkg.QueryerExecer) error {
			atomic.AddInt64(&removeCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var sendCount int64
	var sendErrorCount int64
	sender.EXPECT().Send(gomock.Any()).DoAndReturn(
		func(serviceEvent *subscription.ServiceEvent) error {
			if d.isEmitSendError {
				if serviceEvent.ID%2 == 0 {
					atomic.AddInt64(&sendErrorCount, 1)

					return errors.New("i can't send")
				}
			}
			atomic.AddInt64(&sendCount, 1)

			return nil
		},
	).AnyTimes()

	retranslator := retranslatorpkg.NewRetranslator(
		ctx,
		&config.Retranslator{
			EventChannelSize: d.eventChannelSize,

			MaxConsumers:      d.maxConsumers,
			ConsumerTimeout:   d.consumerTimeout,
			ConsumerBatchTime: d.consumerBatchTime,
			ConsumerBatchSize: d.consumerBatchSize,

			MaxProducers:       d.maxProducers,
			ProducerTimeout:    d.producerTimeout,
			ProducerMaxWorkers: d.producerMaxWorkers,
		},
		repo,
		sender,
	)

	retranslator.Start(ctx)

	time.Sleep(time.Millisecond * 500)

	cancelCtx()

	retranslator.StopWait()

	if !d.isEmitSendError {
		assert.LessOrEqual(t, sendCount, lockCount)
		assert.Equal(t, sendErrorCount, unlockCount)

		return
	}

	assert.Equal(t, sendCount, removeCount)
	assert.Equal(t, lockCount, unlockCount+removeCount)
}

func TestAllEventsCompleteWhenStoppingByContextOne(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			eventChannelSize: 1,

			maxConsumers:      10,
			consumerTimeout:   time.Millisecond,
			consumerBatchTime: time.Microsecond,
			consumerBatchSize: 10,

			maxProducers:       10,
			producerTimeout:    time.Millisecond,
			producerMaxWorkers: 2,

			isEmitSendError: false,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContextTwo(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			eventChannelSize: 10,

			maxConsumers:      20,
			consumerTimeout:   time.Millisecond,
			consumerBatchTime: time.Microsecond,
			consumerBatchSize: 10,

			maxProducers:       20,
			producerTimeout:    time.Millisecond,
			producerMaxWorkers: 2,

			isEmitSendError: true,
		},
	)
}
