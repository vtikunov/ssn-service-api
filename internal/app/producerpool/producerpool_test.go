package producerpool_test

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/app/channellocator"

	"github.com/ozonmp/ssn-service-api/internal/app/producerpool"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/ssn-service-api/internal/app/producer"
	"github.com/ozonmp/ssn-service-api/internal/mocks"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/stretchr/testify/assert"
)

type initData struct {
	numEvents       uint64
	maxProducers    uint64
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
				if serviceEvent.ID%2 == 0 {
					return errors.New("i can't send")
				}
			}
			atomic.AddInt64(&sendCount, 1)

			return nil
		},
	).AnyTimes()

	eventsChannel := make(chan []subscription.ServiceEvent)
	channelLocator := channellocator.NewChannelLocator(eventsChannel)

	producerPool := producerpool.NewProducerPool(d.maxProducers, producer.NewProducerFactory(channelLocator, sender, repo, d.maxWorkers), time.Second)
	producerPool.Start(ctx)

	for i := uint64(1); i <= d.numEvents; i++ {
		eventsChannel <- []subscription.ServiceEvent{{ID: i}}
	}

	producerPool.StopWait()

	assert.LessOrEqual(t, sendCount, int64(d.numEvents))
	assert.Equal(t, int64(d.numEvents), unlockCount+removeCount)
	assert.Equal(t, removeCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByFunc100Events10Producers2Workers(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			numEvents:    100,
			maxProducers: 10,
			maxWorkers:   2,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc1000Events50Producers4Workers(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			numEvents:    1000,
			maxProducers: 50,
			maxWorkers:   4,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByFunc2000Events20Producers3WorkersAndSendError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByFunc(
		t,
		initData{
			numEvents:       2000,
			maxProducers:    20,
			maxWorkers:      3,
			isEmitSendError: true,
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
				if serviceEvent.ID%2 == 0 {
					return errors.New("i can't send")
				}
			}
			atomic.AddInt64(&sendCount, 1)

			return nil
		},
	).AnyTimes()

	eventsChannel := make(chan []subscription.ServiceEvent)
	channelLocator := channellocator.NewChannelLocator(eventsChannel)

	producerPool := producerpool.NewProducerPool(d.maxProducers, producer.NewProducerFactory(channelLocator, sender, repo, d.maxWorkers), time.Second)
	doneChannel := producerPool.Start(ctx)

	for i := uint64(1); i <= d.numEvents; i++ {
		eventsChannel <- []subscription.ServiceEvent{{ID: i}}
	}

	cancelCtx()

	<-doneChannel

	assert.LessOrEqual(t, sendCount, int64(d.numEvents))
	assert.Equal(t, int64(d.numEvents), unlockCount+removeCount)
	assert.Equal(t, removeCount, sendCount)
}

func TestAllEventsCompleteWhenStoppingByContext100Events10Producers2Workers(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			numEvents:    100,
			maxProducers: 10,
			maxWorkers:   2,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext1000Events50Producers4Workers(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			numEvents:    1000,
			maxProducers: 50,
			maxWorkers:   4,
		},
	)
}

func TestAllEventsCompleteWhenStoppingByContext2000Events20Producers3WorkersAndSendError(t *testing.T) {
	SuiteAllEventsCompleteWhenStoppingByContext(
		t,
		initData{
			numEvents:       2000,
			maxProducers:    20,
			maxWorkers:      3,
			isEmitSendError: true,
		},
	)
}
