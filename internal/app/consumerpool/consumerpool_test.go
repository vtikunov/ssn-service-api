package consumerpool_test

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"sync"
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

	eventsChannel := make(chan []subscription.ServiceEvent)

	consumerPool := consumerpool.NewConsumerPool(
		d.maxConsumers,
		consumer.NewConsumerFactory(time.Microsecond, d.batchSize, eventsChannel, repo),
		time.Millisecond,
	)

	var sendCount int64
	doneChannelRoutine := make(chan interface{})
	wg := &sync.WaitGroup{}
	for i := uint64(0); i < d.maxConsumers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case events := <-eventsChannel:
					atomic.AddInt64(&sendCount, int64(len(events)))
				case <-doneChannelRoutine:
					return
				}
			}
		}()
	}

	consumerPool.Start(ctx)

	time.Sleep(time.Millisecond * 500)

	consumerPool.StopWait()

	close(doneChannelRoutine)

	wg.Wait()

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

func TestAllEventsCompleteWhenStoppingByFunc20Consumers50EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			maxConsumers: 20,
			batchSize:    50,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc30Consumers50EventInBatchWithLockError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			maxConsumers:    30,
			batchSize:       50,
			isEmitLockError: true,
		},
	)
}

func SuiteAllEventsCompleteWhenStoppingByContext(t *testing.T, d initData) {
	log.SetOutput(ioutil.Discard)

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

	eventsChannel := make(chan []subscription.ServiceEvent)

	consumerPool := consumerpool.NewConsumerPool(
		d.maxConsumers,
		consumer.NewConsumerFactory(time.Microsecond, d.batchSize, eventsChannel, repo),
		time.Millisecond,
	)

	var sendCount int64
	doneChannelRoutine := make(chan interface{})
	wg := &sync.WaitGroup{}
	for i := uint64(0); i < d.maxConsumers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case events := <-eventsChannel:
					atomic.AddInt64(&sendCount, int64(len(events)))
				case <-doneChannelRoutine:
					return
				}
			}
		}()
	}

	doneChannel := consumerPool.Start(ctx)

	time.Sleep(time.Millisecond * 500)

	cancelCtx()

	<-doneChannel

	close(doneChannelRoutine)

	wg.Wait()

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

func TestAllEventsCompleteWhenStoppingByContext20Consumers50EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			maxConsumers: 20,
			batchSize:    50,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext20Consumers50EventInBatchWithLockError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			maxConsumers:    20,
			batchSize:       50,
			isEmitLockError: true,
		},
	)
}
