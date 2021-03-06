package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/repo"
)

// Consumer - воркер-консьюмер
type Consumer interface {
	Start(ctx context.Context) (doneChannel <-chan interface{})
	StopWait()
}

type eventRepoLockUnlocker interface {
	Lock(ctx context.Context, n uint64, tx repo.QueryerExecer) ([]subscription.ServiceEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64, tx repo.QueryerExecer) error
}

type channelLocator interface {
	GetMainEventsWriteChannel() chan<- []subscription.ServiceEvent
}

type transactionalSession interface {
	Execute(ctx context.Context, fn func(ctx context.Context, tx repo.QueryerExecer) error) error
}

type consumer struct {
	batchTime      time.Duration
	batchSize      uint64
	channelLocator channelLocator
	eventRepo      eventRepoLockUnlocker
	txs            transactionalSession
	doneChannel    chan interface{}
	stopChannel    chan interface{}
	onceStart      *sync.Once
	onceStop       *sync.Once
}

// NewConsumer создает нового воркера-консьюмера.
//
// batchTime: определяет время между обращениями к репозиторию событий.
//
// batchSize: определяет размер пакета событий, получаемый консьюмером из репозитория
// за одно обращение.
//
// channelLocator: локатор каналов.
//
// eventRepo: указатель на экземпляр репозитория событий, из которого консьюмер
// получает и блокирует события.
func NewConsumer(
	ctx context.Context,
	batchTime time.Duration,
	batchSize uint64,
	channelLocator channelLocator,
	eventRepo eventRepoLockUnlocker,
	txs transactionalSession,
) *consumer {

	if batchSize == 0 {
		logger.FatalKV(ctx, "batchSize must be greater than 0")
	}

	return &consumer{
		batchTime:      batchTime,
		batchSize:      batchSize,
		channelLocator: channelLocator,
		eventRepo:      eventRepo,
		txs:            txs,
		onceStart:      &sync.Once{},
		onceStop:       &sync.Once{},
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
					var events []subscription.ServiceEvent
					err := c.txs.Execute(ctx, func(ctx context.Context, tx repo.QueryerExecer) error {
						evs, err := c.eventRepo.Lock(ctx, c.batchSize, tx)
						if err != nil {
							return err
						}

						events = evs

						return nil
					})

					if err != nil {
						logger.ErrorKV(ctx, "consumer: failed to lock events", "err", err)

						continue
					}

					c.channelLocator.GetMainEventsWriteChannel() <- events
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
	batchTime      time.Duration
	batchSize      uint64
	channelLocator channelLocator
	eventRepo      eventRepoLockUnlocker
	txs            transactionalSession
}

// NewConsumerFactory создает фабрику воркеров-консьюмеров.
//
// batchTime: определяет время между обращениями к репозиторию событий.
//
// batchSize: определяет размер пакета событий, получаемый консьюмером из репозитория
// за одно обращение.
//
// channelLocator: локатор каналов.
//
// eventRepo: указатель на экземпляр репозитория событий, из которого консьюмер
// получает и блокирует события.
func NewConsumerFactory(
	batchTime time.Duration,
	batchSize uint64,
	channelLocator channelLocator,
	eventRepo eventRepoLockUnlocker,
	txs transactionalSession,
) *consumerFactory {

	return &consumerFactory{
		batchTime:      batchTime,
		batchSize:      batchSize,
		channelLocator: channelLocator,
		eventRepo:      eventRepo,
		txs:            txs,
	}
}

// Create создает воркера-консьюмера.
func (cf *consumerFactory) Create(ctx context.Context) Consumer {
	return NewConsumer(ctx, cf.batchTime, cf.batchSize, cf.channelLocator, cf.eventRepo, cf.txs)
}
