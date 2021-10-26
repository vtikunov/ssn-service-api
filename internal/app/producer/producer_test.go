package producer_test

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	producerpkg "github.com/ozonmp/ssn-service-api/internal/app/producer"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/ssn-service-api/internal/mocks"
	"github.com/stretchr/testify/assert"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

type initData struct {
	numEvents       uint64
	maxWorkers      uint64
	isEmitSendError bool
}

func SuiteAllEventsCompleteWhenStoppingByFunc(t *testing.T, d initData) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	ctrl := gomock.NewController(t)
	ctx := context.Background()
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

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
	sender.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, serviceEvent *subscription.ServiceEvent) error {
			if d.isEmitSendError {
				count := atomic.LoadInt64(&sendCount)
				if count%2 == 0 {
					return errors.New("i can't send")
				}
			}
			atomic.AddInt64(&sendCount, 1)

			return nil
		},
	).AnyTimes()

	eventsChannel := make(chan subscription.ServiceEvent)

	producer := producerpkg.NewProducer(time.Second, eventsChannel, sender, repo, d.maxWorkers)

	producer.Start(ctx)

	for i := uint64(1); i <= d.numEvents; i++ {
		eventsChannel <- subscription.ServiceEvent{ID: i}
	}

	producer.StopWait()

	assert.LessOrEqual(t, sendCount, int64(d.numEvents))
	assert.Equal(t, int64(d.numEvents), unlockCount+removeCount)
	assert.Equal(t, removeCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByFunc100Events1Worker(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			numEvents:  100,
			maxWorkers: 1,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc1000Events2Worker(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			numEvents:  1000,
			maxWorkers: 2,
		},
	)
}

func SuiteAllEventsCompleteWhenStoppingByContext(t *testing.T, d initData) {
	log.SetOutput(ioutil.Discard)

	t.Parallel()

	ctrl := gomock.NewController(t)
	ctx, cancelCtx := context.WithCancel(context.Background())
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

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
	sender.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, serviceEvent *subscription.ServiceEvent) error {
			if d.isEmitSendError {
				count := atomic.LoadInt64(&sendCount)
				if count%2 == 0 {
					return errors.New("i can't send")
				}
			}
			atomic.AddInt64(&sendCount, 1)

			return nil
		},
	).AnyTimes()

	eventsChannel := make(chan subscription.ServiceEvent)

	producer := producerpkg.NewProducer(time.Second, eventsChannel, sender, repo, d.maxWorkers)

	doneChannel := producer.Start(ctx)

	for i := uint64(1); i <= d.numEvents; i++ {
		eventsChannel <- subscription.ServiceEvent{ID: i}
	}

	cancelCtx()

	<-doneChannel

	assert.LessOrEqual(t, sendCount, int64(d.numEvents))
	assert.Equal(t, int64(d.numEvents), unlockCount+removeCount)
	assert.Equal(t, removeCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByContext100Events1Worker(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			numEvents:  100,
			maxWorkers: 1,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext1000Events2Worker(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			numEvents:  1000,
			maxWorkers: 2,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext1000Events2WorkerEndSendError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			numEvents:       1000,
			maxWorkers:      2,
			isEmitSendError: true,
		},
	)
}
