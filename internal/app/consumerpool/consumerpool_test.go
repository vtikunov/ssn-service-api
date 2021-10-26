package consumerpool_test

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/ssn-service-api/internal/app/consumer"
	"github.com/ozonmp/ssn-service-api/internal/app/consumerpool"
	"github.com/ozonmp/ssn-service-api/internal/mocks"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/stretchr/testify/assert"
)

type initData struct {
	maxConsumers    uint64
	batchSize       uint64
	isEmitLockError bool
}

func SuiteAllEventsCompleteWhenStoppingByFunc(t *testing.T, d initData) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	ctrl := gomock.NewController(t)
	ctx := context.Background()
	repo := mocks.NewMockEventRepo(ctrl)

	var lockCount int64
	var lockVoidCount int64
	repo.EXPECT().Lock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error) {
			defer atomic.AddInt64(&lockVoidCount, 1)

			if d.isEmitLockError && atomic.LoadInt64(&lockVoidCount)%2 > 0 {
				return nil, errors.New("i can't lock")
			}

			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i}
			}

			atomic.AddInt64(&lockCount, int64(n))

			return result, nil
		},
	).AnyTimes()

	eventsChannel := make(chan subscription.ServiceEvent)

	consumerPool := consumerpool.NewConsumerPool(
		d.maxConsumers,
		consumer.NewConsumerFactory(time.Microsecond, d.batchSize, eventsChannel, repo),
		time.Millisecond,
	)

	consumerPool.Start(ctx)

	var sendCount int64
	doneChannelRoutine := make(chan interface{})
	go func() {
		for {
			select {
			case <-eventsChannel:
				atomic.AddInt64(&sendCount, 1)
			case <-doneChannelRoutine:
				return
			}
		}
	}()

	time.Sleep(time.Millisecond * 500)

	consumerPool.StopWait()

	close(doneChannelRoutine)

	assert.Equal(t, lockCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByFunc10Consumers1EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			maxConsumers: 10,
			batchSize:    1,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc100Consumers500EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			maxConsumers: 100,
			batchSize:    500,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc100Consumers500EventInBatchWithLockError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			maxConsumers:    100,
			batchSize:       500,
			isEmitLockError: true,
		},
	)
}

func SuiteAllEventsCompleteWhenStoppingByContext(t *testing.T, d initData) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	ctrl := gomock.NewController(t)
	ctx, cancelCtx := context.WithCancel(context.Background())
	repo := mocks.NewMockEventRepo(ctrl)

	var lockCount int64
	var lockVoidCount int64
	repo.EXPECT().Lock(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error) {
			defer atomic.AddInt64(&lockVoidCount, 1)

			if d.isEmitLockError && atomic.LoadInt64(&lockVoidCount)%2 > 0 {
				return nil, errors.New("i can't lock")
			}

			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i}
			}

			atomic.AddInt64(&lockCount, int64(n))

			return result, nil
		},
	).AnyTimes()

	eventsChannel := make(chan subscription.ServiceEvent)

	consumerPool := consumerpool.NewConsumerPool(
		d.maxConsumers,
		consumer.NewConsumerFactory(time.Microsecond, d.batchSize, eventsChannel, repo),
		time.Millisecond,
	)

	doneChannel := consumerPool.Start(ctx)

	time.Sleep(time.Millisecond * 500)

	cancelCtx()

	var sendCount int64
	doneChannelRoutine := make(chan interface{})
	go func() {
		for {
			select {
			case <-eventsChannel:
				atomic.AddInt64(&sendCount, 1)
			case <-doneChannelRoutine:
				return
			}
		}
	}()

	<-doneChannel

	close(doneChannelRoutine)

	assert.Equal(t, lockCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByContext10Consumers1EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			maxConsumers: 10,
			batchSize:    1,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext100Consumers500EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			maxConsumers: 100,
			batchSize:    500,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext100Consumers500EventInBatchWithLockError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			maxConsumers:    100,
			batchSize:       500,
			isEmitLockError: true,
		},
	)
}
