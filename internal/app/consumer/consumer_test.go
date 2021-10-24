package consumer_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	consumerpkg "github.com/ozonmp/ssn-service-api/internal/app/consumer"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/ssn-service-api/internal/mocks"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/stretchr/testify/assert"
)

func TestAllEventsComplete(t *testing.T) {
	t.Parallel()
	log.SetOutput(ioutil.Discard)

	tests := []struct {
		name            string
		batchSize       uint64
		isEmitLockError bool
	}{
		{
			batchSize: 1,
		},
		{
			batchSize: 100,
		},
		{
			batchSize: 500,
		},
		{
			batchSize:       1,
			isEmitLockError: true,
		},
		{
			batchSize:       1000,
			isEmitLockError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf(
				"Lock %v events in batch. Is emit lock errors: %v",
				tt.batchSize,
				tt.isEmitLockError,
			),
			func(t *testing.T) {
				t.Parallel()

				ctrl := gomock.NewController(t)
				repo := mocks.NewMockEventRepo(ctrl)

				var lockCount int64
				var lockVoidCount int64
				repo.EXPECT().Lock(gomock.Any()).DoAndReturn(
					func(n uint64) ([]subscription.ServiceEvent, error) {
						defer atomic.AddInt64(&lockVoidCount, 1)

						if tt.isEmitLockError && atomic.LoadInt64(&lockVoidCount)%2 > 0 {
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

				consumer := consumerpkg.NewConsumer(time.Microsecond, tt.batchSize, eventsChannel, repo)

				doneChannel := consumer.Start()

				var sendCount int64
				doneChannelRoutine := make(chan interface{})
				go func() {
					defer close(doneChannelRoutine)
					for {
						select {
						case <-eventsChannel:
							atomic.AddInt64(&sendCount, 1)
						case <-doneChannel:
							return
						}
					}
				}()

				time.Sleep(time.Millisecond * 500)

				consumer.StopWait()

				<-doneChannelRoutine

				assert.Greater(t, lockCount, int64(0))
				assert.Equal(t, lockCount, sendCount)
			},
		)
	}
}
