package retranslator_test

import (
	"context"
	"errors"
	"fmt"
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

func TestAllEventsComplete(t *testing.T) {
	t.Parallel()
	log.SetOutput(ioutil.Discard)

	tests := []struct {
		eventChannelSize uint64

		maxConsumers      uint64
		consumerTimeout   time.Duration
		consumerBatchTime time.Duration
		consumerBatchSize uint64

		maxProducers       uint64
		producerTimeout    time.Duration
		producerMaxWorkers uint64

		isEmitSendError bool
	}{
		{
			eventChannelSize:   1,
			maxConsumers:       1,
			consumerTimeout:    time.Millisecond,
			consumerBatchTime:  time.Microsecond,
			consumerBatchSize:  10,
			maxProducers:       1,
			producerTimeout:    time.Millisecond,
			producerMaxWorkers: 2,
		},
		{
			eventChannelSize:   50,
			maxConsumers:       10,
			consumerTimeout:    time.Millisecond,
			consumerBatchTime:  time.Microsecond,
			consumerBatchSize:  100,
			maxProducers:       10,
			producerTimeout:    time.Millisecond,
			producerMaxWorkers: 2,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf(
				"Retranslate with %v consumers and %v producers. Is emit send errors: %v",
				tt.maxConsumers,
				tt.maxProducers,
				tt.isEmitSendError,
			),
			func(t *testing.T) {
				t.Parallel()

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
				repo.EXPECT().Unlock(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, eventIDs []uint64) error {
						atomic.AddInt64(&unlockCount, int64(len(eventIDs)))

						return nil
					},
				).AnyTimes()

				var removeCount int64
				repo.EXPECT().Remove(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, eventIDs []uint64) error {
						atomic.AddInt64(&removeCount, int64(len(eventIDs)))

						return nil
					},
				).AnyTimes()

				var sendCount int64
				var sendErrorCount int64
				sender.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, serviceEvent *subscription.ServiceEvent) error {
						if tt.isEmitSendError {
							count := atomic.LoadInt64(&sendCount)
							if count%2 == 0 {
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
						EventChannelSize: tt.eventChannelSize,
						EventRepo:        repo,
						EventSender:      sender,

						MaxConsumers:      tt.maxConsumers,
						ConsumerTimeout:   tt.consumerTimeout,
						ConsumerBatchTime: tt.consumerBatchTime,
						ConsumerBatchSize: tt.consumerBatchSize,

						MaxProducers:       tt.maxProducers,
						ProducerTimeout:    tt.producerTimeout,
						ProducerMaxWorkers: tt.producerMaxWorkers,
					},
				)

				retranslator.Start(ctx)

				time.Sleep(time.Millisecond * 500)

				retranslator.StopWait()

				if !tt.isEmitSendError {
					assert.LessOrEqual(t, sendCount, lockCount)
					assert.Equal(t, sendErrorCount, unlockCount)

					return
				}

				assert.Equal(t, sendCount, lockCount)
				assert.Equal(t, sendCount, removeCount)
				assert.Equal(t, lockCount, unlockCount+removeCount)
			},
		)
	}
}
