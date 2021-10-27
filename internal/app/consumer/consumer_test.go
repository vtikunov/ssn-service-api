package consumer_test

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/app/channellocator"

	"github.com/golang/mock/gomock"
	consumerpkg "github.com/ozonmp/ssn-service-api/internal/app/consumer"
	"github.com/ozonmp/ssn-service-api/internal/mocks"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/stretchr/testify/assert"
)

type initData struct {
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
	repo.EXPECT().LockExceptLockedByServiceID(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error) {
			defer atomic.AddInt64(&lockVoidCount, 1)

			if d.isEmitLockError && atomic.LoadInt64(&lockVoidCount)%2 > 0 {
				return nil, errors.New("i can't lock")
			}

			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i, Service: &subscription.Service{ID: start + i}}
			}

			return result, nil
		},
	).AnyTimes()

	repo.EXPECT().LockByServiceID(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, serviceID uint64) ([]subscription.ServiceEvent, error) {
			atomic.AddInt64(&lockCount, 1)

			return []subscription.ServiceEvent{{ID: serviceID, Service: &subscription.Service{ID: serviceID}}}, nil
		},
	).AnyTimes()

	var unlockCount int64
	repo.EXPECT().Unlock(gomock.Any()).DoAndReturn(
		func(eventIDs []uint64) error {
			atomic.AddInt64(&unlockCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	eventsChannel := make(chan []subscription.ServiceEvent)
	channelLocator := channellocator.NewChannelLocator(eventsChannel)

	consumer := consumerpkg.NewConsumer(time.Microsecond, d.batchSize, channelLocator, repo)

	doneChannel := consumer.Start(ctx)

	var sendCount int64
	doneChannelRoutine := make(chan interface{})
	go func() {
		defer close(doneChannelRoutine)
		for {
			select {
			case events := <-eventsChannel:
				for _, event := range events {
					endEventsChannel, err := channelLocator.GetEventsServiceIDReadChannel(event.Service.ID)
					if err != nil {
						continue
					}
					endEvents := <-endEventsChannel
					atomic.AddInt64(&sendCount, int64(len(endEvents)))
				}
			case <-doneChannel:
				return
			}
		}
	}()

	time.Sleep(time.Millisecond * 500)

	consumer.StopWait()

	<-doneChannelRoutine

	assert.Equal(t, lockCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByFunc1EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			batchSize: 1,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc100EventsInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			batchSize: 100,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc200EventsInBatchEndLockError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			batchSize:       200,
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
	repo.EXPECT().LockExceptLockedByServiceID(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error) {
			defer atomic.AddInt64(&lockVoidCount, 1)

			if d.isEmitLockError && atomic.LoadInt64(&lockVoidCount)%2 > 0 {
				return nil, errors.New("i can't lock")
			}

			start := uint64(atomic.LoadInt64(&lockCount)) + 1
			result := make([]subscription.ServiceEvent, n)

			for i := uint64(0); i < n; i++ {
				result[i] = subscription.ServiceEvent{ID: start + i, Service: &subscription.Service{ID: start + i}}
			}

			return result, nil
		},
	).AnyTimes()

	repo.EXPECT().LockByServiceID(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, serviceID uint64) ([]subscription.ServiceEvent, error) {
			atomic.AddInt64(&lockCount, 1)

			return []subscription.ServiceEvent{{ID: serviceID, Service: &subscription.Service{ID: serviceID}}}, nil
		},
	).AnyTimes()

	var unlockCount int64
	repo.EXPECT().Unlock(gomock.Any()).DoAndReturn(
		func(eventIDs []uint64) error {
			atomic.AddInt64(&unlockCount, int64(len(eventIDs)))

			return nil
		},
	).AnyTimes()

	eventsChannel := make(chan []subscription.ServiceEvent)
	channelLocator := channellocator.NewChannelLocator(eventsChannel)

	consumer := consumerpkg.NewConsumer(time.Microsecond, d.batchSize, channelLocator, repo)

	doneChannel := consumer.Start(ctx)

	var sendCount int64
	doneChannelRoutine := make(chan interface{})
	go func() {
		defer close(doneChannelRoutine)
		for {
			select {
			case events := <-eventsChannel:
				for _, event := range events {
					endEventsChannel, err := channelLocator.GetEventsServiceIDReadChannel(event.Service.ID)
					if err != nil {
						continue
					}
					endEvents := <-endEventsChannel
					atomic.AddInt64(&sendCount, int64(len(endEvents)))
				}
			case <-doneChannel:
				return
			}
		}
	}()

	time.Sleep(time.Millisecond * 500)

	cancelCtx()

	<-doneChannelRoutine

	assert.Equal(t, lockCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByContext1EventInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			batchSize: 1,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext100EventsInBatch(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			batchSize: 100,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext200EventsInBatchEndLockError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			batchSize:       200,
			isEmitLockError: true,
		},
	)
}
