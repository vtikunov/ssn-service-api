package consumerpool

import (
	"log"
	"sync/atomic"
	"time"

	consumerpkg "github.com/ozonmp/ssn-service-api/internal/app/consumer"
)

type consumerFactory interface {
	Create() consumerpkg.Consumer
}

type consumerPool struct {
	maxConsumers    int64
	consumerFactory consumerFactory
	consumerTimeout time.Duration
	doneChannel     chan interface{}
	stopChannel     chan interface{}
	isStarted       bool
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
	maxConsumers uint64,
	consumerFactory consumerFactory,
	consumerTimeout time.Duration,
) *consumerPool {

	if maxConsumers == 0 {
		maxConsumers = 1
	}

	return &consumerPool{
		maxConsumers:    int64(maxConsumers),
		consumerFactory: consumerFactory,
		consumerTimeout: consumerTimeout,
	}
}

// Start запускает работу пула.
func (cp *consumerPool) Start() {
	if cp.isStarted {
		log.Panic("pull is already started")
	}
	cp.isStarted = true
	cp.doneChannel = make(chan interface{})
	cp.stopChannel = make(chan interface{})

	go cp.dispatch()
}

func (cp *consumerPool) dispatch() {
	defer close(cp.doneChannel)

	var consumerCount int64

	for {
		select {
		case <-cp.stopChannel:
			if atomic.LoadInt64(&consumerCount) == 0 {
				return
			}
		default:
			if atomic.LoadInt64(&consumerCount) >= cp.maxConsumers {
				break
			}

			atomic.AddInt64(&consumerCount, 1)
			go func() {
				defer atomic.AddInt64(&consumerCount, -1)
				consumer := cp.consumerFactory.Create()
				doneChannel := consumer.Start()
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

// StopWait отправляет команду Stop всем работающим воркерам,
// дожидается окончания их работы и останавливает работу пула.
//
// Обратите внимание! Метод возвращает return после остановки пула.
func (cp *consumerPool) StopWait() {
	if !cp.isStarted {
		log.Panic("pull is already stopped")
	}
	close(cp.stopChannel)
	<-cp.doneChannel

	cp.isStarted = false
}
