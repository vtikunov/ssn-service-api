package consumerpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	consumerpkg "github.com/ozonmp/ssn-service-api/internal/retranslator/consumer"
)

type consumerFactory interface {
	Create(ctx context.Context) consumerpkg.Consumer
}

type consumerPool struct {
	maxConsumers    int64
	consumerFactory consumerFactory
	consumerTimeout time.Duration
	doneChannel     chan interface{}
	stopChannel     chan interface{}
	onceStart       *sync.Once
	onceStop        *sync.Once
}

// NewConsumerPool создает пул воркеров-консьюмеров.
//
// maxConsumers: определяет максимальное количество воркеров-консьюмеров,
// которые будут запущены конкуретно, при окончании работы какого-либо воркера-консьюмера
// пул будет создавать и запускать следующий в пределах указанного лимита.
//
// consumerFactory: принимает фабрику, необходимую для создания
// каждого нового экземпляра воркера-консьюмера.
//
// consumerTimeout: определяет максимальное время работы каждого нового
// экземпляра воркера-консьюемра, по истечении которого ему будет направлена команда Stop.
func NewConsumerPool(
	ctx context.Context,
	maxConsumers uint64,
	consumerFactory consumerFactory,
	consumerTimeout time.Duration,
) *consumerPool {

	if maxConsumers == 0 {
		logger.FatalKV(ctx, "maxConsumers must be greater than 0")
	}

	return &consumerPool{
		maxConsumers:    int64(maxConsumers),
		consumerFactory: consumerFactory,
		consumerTimeout: consumerTimeout,
		onceStart:       &sync.Once{},
		onceStop:        &sync.Once{},
	}
}

// Start запускает работу пула.
func (cp *consumerPool) Start(ctx context.Context) (doneChannel <-chan interface{}) {
	cp.onceStart.Do(func() {
		cp.doneChannel = make(chan interface{})
		cp.stopChannel = make(chan interface{})

		go cp.dispatch(ctx)
	})

	return cp.doneChannel
}

func (cp *consumerPool) dispatch(ctx context.Context) {
	defer close(cp.doneChannel)

	var consumerCount int64

	for {
		select {
		case <-cp.stopChannel:
			if atomic.LoadInt64(&consumerCount) == 0 {
				return
			}
		case <-ctx.Done():
			cp.stop()
		default:
			if atomic.LoadInt64(&consumerCount) >= cp.maxConsumers {
				break
			}

			atomic.AddInt64(&consumerCount, 1)
			go func() {
				defer atomic.AddInt64(&consumerCount, -1)
				consumer := cp.consumerFactory.Create(ctx)
				doneChannel := consumer.Start(ctx)
				timeout := time.NewTimer(cp.consumerTimeout)

				select {
				case <-doneChannel:
				case <-cp.stopChannel:
					consumer.StopWait()
				case <-timeout.C:
					consumer.StopWait()
				}
			}()
		}
	}
}

func (cp *consumerPool) stop() {
	cp.onceStop.Do(func() {
		close(cp.stopChannel)
	})
}

// StopWait отправляет команду Stop всем работающим воркерам,
// дожидается окончания их работы и останавливает работу пула.
//
// Обратите внимание! Метод возвращает return после остановки пула.
func (cp *consumerPool) StopWait() {
	cp.stop()
	<-cp.doneChannel
}
