package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

type Consumer interface {
	Start(ctx context.Context) (doneChannel <-chan interface{})
	StopWait()
}

type eventRepoLocker interface {
	Lock(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error)
}

type consumer struct {
	batchTime     time.Duration
	batchSize     uint64
	eventsChannel chan<- subscription.ServiceEvent
	eventRepo     eventRepoLocker
	doneChannel   chan interface{}
	stopChannel   chan interface{}
	onceStart     *sync.Once
	onceStop      *sync.Once
}

// NewConsumer создает нового воркера-консьюмера.
//
// batchTime: определяет время между обращениями к репозиторию событий.
//
// batchSize: определяет размер пакета событий, получаемый консьюмером из репозитория
// за одно обращение.
//
// eventsChannel: канал для записи событий консьюмером.
//
// eventRepo: указатель на экземпляр репозитория событий, из которого консьюмер
// получает и блокирует события.
func NewConsumer(
	batchTime time.Duration,
	batchSize uint64,
	eventsChannel chan<- subscription.ServiceEvent,
	eventRepo eventRepoLocker,
) *consumer {

	if batchSize == 0 {
		log.Panicln("batchSize must be greater than 0")
	}

	return &consumer{
		batchTime:     batchTime,
		batchSize:     batchSize,
		eventsChannel: eventsChannel,
		eventRepo:     eventRepo,
		onceStart:     &sync.Once{},
		onceStop:      &sync.Once{},
	}
}

// Start запускает работу консьюмера.
//
// Возвращает канал для чтения doneChannel, который закрывается консьюмером при его остановке.
func (c *consumer) Start(ctx context.Context) (doneChannel <-chan interface{}) {
	c.onceStart.Do(func() {
		c.doneChannel = make(chan interface{})
		c.stopChannel = make(chan interface{})

		go func() {
			defer close(c.doneChannel)
			timeout := time.NewTimer(c.batchTime)

			for {
				select {
				case <-ctx.Done():
					c.stop()

					return
				case <-c.stopChannel:
					return
				case <-timeout.C:
					events, err := c.eventRepo.Lock(ctx, c.batchSize)
					if err != nil {
						log.Printf("consumer: failed to lock events - %v", err)

						continue
					}
					for _, event := range events {
						c.eventsChannel <- event
					}
				}
			}
		}()
	})

	return c.doneChannel
}

func (c *consumer) stop() {
	c.onceStop.Do(func() {
		close(c.stopChannel)
	})
}

// StopWait отправляет команду Stop консьюмеру,
// дожидается окончания его работы и останавливает работу консьюмера.
//
// Обратите внимание! Метод возвращает return после остановки консьюмера.
func (c *consumer) StopWait() {
	c.stop()
	<-c.doneChannel
}

type consumerFactory struct {
	batchTime     time.Duration
	batchSize     uint64
	eventsChannel chan<- subscription.ServiceEvent
	eventRepo     eventRepoLocker
}

// NewConsumerFactory создает фабрику воркеров-консьюмеров.
//
// batchTime: определяет время между обращениями к репозиторию событий.
//
// batchSize: определяет размер пакета событий, получаемый консьюмером из репозитория
// за одно обращение.
//
// eventsChannel: канал для записи событий консьюмером.
//
// eventRepo: указатель на экземпляр репозитория событий, из которого консьюмер
// получает и блокирует события.
func NewConsumerFactory(
	batchTime time.Duration,
	batchSize uint64,
	eventsChannel chan<- subscription.ServiceEvent,
	eventRepo eventRepoLocker,
) *consumerFactory {

	return &consumerFactory{
		batchTime:     batchTime,
		batchSize:     batchSize,
		eventsChannel: eventsChannel,
		eventRepo:     eventRepo,
	}
}

// Create создает воркера-консьюмера.
func (cf *consumerFactory) Create() Consumer {
	return NewConsumer(cf.batchTime, cf.batchSize, cf.eventsChannel, cf.eventRepo)
}
