package retranslator_test

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	retranslatorpkg "github.com/ozonmp/ssn-service-api/internal/app/retranslator"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/ssn-service-api/internal/mocks"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/stretchr/testify/assert"
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

func SuiteAllEventsCompleteWhenStoppingByFunc(t *testing.T, d initData) {
	log.SetOutput(ioutil.Discard)

	ctrl := gomock.NewController(t)
	ctx := context.Background()
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	var lockCount int64
	repo.EXPECT().Lock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error) {
			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i}
			}

			atomic.AddInt64(&lockCount, int64(n))

			return result, nil
		},
	).AnyTimes()

	var unlockCount int64
	repo.EXPECT().Unlock(gomock.Any()).DoAndReturn(
		func(eventIDs []uint64) error {
			atomic.AddInt64(&unlockCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var removeCount int64
	repo.EXPECT().Remove(gomock.Any()).DoAndReturn(
		func(eventIDs []uint64) error {
			atomic.AddInt64(&removeCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var sendCount int64
	var sendErrorCount int64
	sender.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, serviceEvent *subscription.ServiceEvent) error {
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
		&retranslatorpkg.Configuration{
			EventChannelSize: d.eventChannelSize,
			EventRepo:        repo,
			EventSender:      sender,

			MaxConsumers:      d.maxConsumers,
			ConsumerTimeout:   d.consumerTimeout,
			ConsumerBatchTime: d.consumerBatchTime,
			ConsumerBatchSize: d.consumerBatchSize,

			MaxProducers:       d.maxProducers,
			ProducerTimeout:    d.producerTimeout,
			ProducerMaxWorkers: d.producerMaxWorkers,
		},
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
	log.SetOutput(ioutil.Discard)

	ctrl := gomock.NewController(t)
	ctx, cancelCtx := context.WithCancel(context.Background())
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	var lockCount int64
	repo.EXPECT().Lock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error) {
			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i}
			}

			atomic.AddInt64(&lockCount, int64(n))

			return result, nil
		},
	).AnyTimes()

	var unlockCount int64
	repo.EXPECT().Unlock(gomock.Any()).DoAndReturn(
		func(eventIDs []uint64) error {
			atomic.AddInt64(&unlockCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var removeCount int64
	repo.EXPECT().Remove(gomock.Any()).DoAndReturn(
		func(eventIDs []uint64) error {
			atomic.AddInt64(&removeCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	var sendCount int64
	var sendErrorCount int64
	sender.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, serviceEvent *subscription.ServiceEvent) error {
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
		&retranslatorpkg.Configuration{
			EventChannelSize: d.eventChannelSize,
			EventRepo:        repo,
			EventSender:      sender,

			MaxConsumers:      d.maxConsumers,
			ConsumerTimeout:   d.consumerTimeout,
			ConsumerBatchTime: d.consumerBatchTime,
			ConsumerBatchSize: d.consumerBatchSize,

			MaxProducers:       d.maxProducers,
			ProducerTimeout:    d.producerTimeout,
			ProducerMaxWorkers: d.producerMaxWorkers,
		},
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
