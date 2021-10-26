package producerpool

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	producerpkg "github.com/ozonmp/ssn-service-api/internal/app/producer"
)

type producerFactory interface {
	Create(timeout time.Duration) producerpkg.Producer
}

type producerPool struct {
	maxProducers    int64
	producerFactory producerFactory
	producerTimeout time.Duration
	doneChannel     chan interface{}
	stopChannel     chan interface{}
	isStarted       bool
}

// NewProducerPool создает пул воркеров-продьюсеров.
//
// maxProducers: определяет максимальное количество воркеров-продьюсеров,
// которые будут запущены конкуретно, при окончании работы какого-либо воркера-продьюсера
// пул будет создавать и запускать следующий в пределах указанного лимита.
//
// producerFactory: принимает фабрику, необходимую для создания
// каждого нового экземпляра воркера-продьюсера.
//
// producerTimeout: определяет максимальное время работы каждого нового
// экземпляра воркера-продьюсера, по истечении которого ему будет направлена команда Stop.
func NewProducerPool(
	maxProducers uint64,
	producerFactory producerFactory,
	producerTimeout time.Duration,
) *producerPool {

	if maxProducers == 0 {
		maxProducers = 1
	}

	return &producerPool{
		maxProducers:    int64(maxProducers),
		producerFactory: producerFactory,
		producerTimeout: producerTimeout,
	}
}

// Start запускает работу пула.
func (pp *producerPool) Start(ctx context.Context) {
	if pp.isStarted {
		log.Panic("pull is already started")
	}
	pp.isStarted = true
	pp.doneChannel = make(chan interface{})
	pp.stopChannel = make(chan interface{})

	go pp.dispatch(ctx)
}

func (pp *producerPool) dispatch(ctx context.Context) {
	defer close(pp.doneChannel)

	var producerCount int64

	for {
		select {
		case <-ctx.Done():
			if atomic.LoadInt64(&producerCount) == 0 {
				return
			}
		case <-pp.stopChannel:
			if atomic.LoadInt64(&producerCount) == 0 {
				return
			}
		default:
			if atomic.LoadInt64(&producerCount) >= pp.maxProducers {
				break
			}

			atomic.AddInt64(&producerCount, 1)
			go func() {
				defer atomic.AddInt64(&producerCount, -1)
				producer := pp.producerFactory.Create(pp.producerTimeout)
				doneChannel := producer.Start(ctx)
				timeout := time.NewTimer(pp.producerTimeout)

				select {
				case <-doneChannel:
				case <-ctx.Done():
					producer.StopWait()
				case <-pp.stopChannel:
					producer.StopWait()
				case <-timeout.C:
					producer.StopWait()
				}
			}()
		}
	}
}

// StopWait отправляет команду Stop всем работающим воркерам,
// дожидается окончания их работы и останавливает работу пула.
//
// Обратите внимание! Метод возвращает return после остановки пула.
func (pp *producerPool) StopWait() {
	if !pp.isStarted {
		log.Panic("pull is already stopped")
	}
	close(pp.stopChannel)
	<-pp.doneChannel

	pp.isStarted = false
}
