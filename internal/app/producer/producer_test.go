package producer_test

import (
	"errors"
	"fmt"
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

func TestAllEventsComplete(t *testing.T) {
	t.Parallel()
	log.SetOutput(ioutil.Discard)

	tests := []struct {
		numEvents       uint64
		maxWorkers      uint64
		isEmitSendError bool
	}{
		{
			numEvents:  100,
			maxWorkers: 1,
		},
		{
			numEvents:  1000,
			maxWorkers: 2,
		},
		{
			numEvents:  10000,
			maxWorkers: 10,
		},
		{
			numEvents:       100,
			maxWorkers:      2,
			isEmitSendError: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf(
				"Send %v events with %v workers. Is emit send errors: %v",
				tt.numEvents,
				tt.maxWorkers,
				tt.isEmitSendError,
			),
			func(t *testing.T) {

				ctrl := gomock.NewController(t)
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
				sender.EXPECT().Send(gomock.Any()).DoAndReturn(
					func(serviceEvent *subscription.ServiceEvent) error {
						if tt.isEmitSendError {
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

				producer := producerpkg.NewProducer(time.Second, eventsChannel, sender, repo, tt.maxWorkers)

				producer.Start()

				for i := uint64(0); i < tt.numEvents; i++ {
					eventsChannel <- subscription.ServiceEvent{ID: i}
				}

				producer.StopWait()

				assert.LessOrEqual(t, sendCount, int64(tt.numEvents))
				assert.Equal(t, int64(tt.numEvents), unlockCount+removeCount)
				assert.Equal(t, removeCount, sendCount)
			},
		)
	}
}
