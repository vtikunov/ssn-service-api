package consumerpool_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/app/consumerpool"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/ssn-service-api/internal/app/consumer"
	"github.com/ozonmp/ssn-service-api/internal/mocks"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/stretchr/testify/assert"
)

func TestAllEventsComplete(t *testing.T) {
	t.Parallel()
	log.SetOutput(ioutil.Discard)

	tests := []struct {
		maxConsumers    uint64
		batchSize       uint64
		isEmitLockError bool
	}{
		{
			maxConsumers: 1,
			batchSize:    1,
		},
		{
			maxConsumers: 2,
			batchSize:    100,
		},
		{
			maxConsumers: 5,
			batchSize:    500,
		},
		{
			maxConsumers:    1,
			batchSize:       1,
			isEmitLockError: true,
		},
		{
			maxConsumers:    10,
			batchSize:       1000,
			isEmitLockError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf(
				"Lock %v events in batch with %v consumers. Is emit lock errors: %v",
				tt.batchSize,
				tt.maxConsumers,
				tt.isEmitLockError,
			),
			func(t *testing.T) {
				t.Parallel()

				ctrl := gomock.NewController(t)
				ctx := context.Background()
				repo := mocks.NewMockEventRepo(ctrl)

				var lockCount int64
				var lockVoidCount int64
				repo.EXPECT().Lock(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error) {
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

				consumerPool := consumerpool.NewConsumerPool(
					tt.maxConsumers,
					consumer.NewConsumerFactory(time.Microsecond, tt.batchSize, eventsChannel, repo),
					time.Millisecond,
				)

				consumerPool.Start(ctx)

				var sendCount int64
				doneChannelRoutine := make(chan interface{})
				exitChannelRoutine := make(chan interface{})
				go func() {
					for {
						select {
						case <-eventsChannel:
							atomic.AddInt64(&sendCount, 1)
						case <-doneChannelRoutine:
							close(exitChannelRoutine)

							return
						}
					}
				}()

				time.Sleep(time.Millisecond * 500)

				consumerPool.StopWait()

				close(doneChannelRoutine)

				<-exitChannelRoutine

				assert.Equal(t, lockCount, sendCount)
			},
		)
	}
}
