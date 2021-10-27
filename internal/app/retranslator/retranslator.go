package retranslator

import (
	"context"
	"sync"
	"time"

	"github.com/ozonmp/ssn-service-api/internal/app/channellocator"

	"github.com/ozonmp/ssn-service-api/internal/app/consumer"
	"github.com/ozonmp/ssn-service-api/internal/app/consumerpool"
	"github.com/ozonmp/ssn-service-api/internal/app/producer"
	"github.com/ozonmp/ssn-service-api/internal/app/producerpool"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

type consumerPool interface {
	Start(ctx context.Context) (doneChannel <-chan interface{})
	StopWait()
}

type producerPool interface {
	Start(ctx context.Context) (doneChannel <-chan interface{})
	StopWait()
}

type eventRepo interface {
	Lock(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error)
	LockExceptLockedByServiceID(ctx context.Context, n uint64) ([]subscription.ServiceEvent, error)
	LockByServiceID(ctx context.Context, serviceID uint64) ([]subscription.ServiceEvent, error)
	Unlock(eventIDs []uint64) error

	Remove(eventIDs []uint64) error
}

type eventSender interface {
	Send(ctx context.Context, serviceEvent *subscription.ServiceEvent) error
}

// Configuration - структура для передачи настроек ретранслятору.
//
// EventChannelSize: размер канала событий между пулами консьюмеров и продьюсеров.
//
// EventRepo: указатель на экземпляр репозитория событий.
//
// EventSender: указатель на экземпляр сендера, куда продьюсер перенаправляет события.
//
// MaxConsumers: определяет максимальное количество воркеров-консьюмеров,
// которые будут запущены конкуретно, при окончании работы какого-либо воркера-консьюмера
// пул консьюмеров будет создавать и запускать следующий в пределах указанного лимита.
//
// ConsumerTimeout: определяет максимальное время работы каждого нового экземпляра воркера-консьюемра,
// по истечении которого ему будет направлена команда Stop.
//
// ConsumerBatchTime: определяет время между обращениями к репозиторию событий.
//
// ConsumerBatchSize: определяет размер пакета событий, получаемый консьюмером из репозитория
// за одно обращение.
//
// MaxProducers: определяет максимальное количество воркеров-продьюсеров,
// которые будут запущены конкуретно, при окончании работы какого-либо воркера-продьюсера
// пул продьюсеров будет создавать и запускать следующий в пределах указанного лимита.
//
// ProducerTimeout: определяет максимальное время работы каждого нового
// экземпляра воркера-продьюсера, по истечении которого ему будет направлена команда Stop.
//
// ProducerMaxWorkers: пределяет максимальное количество вспомогательных воркеров работы продьюсера
// с репозиторием событий eventRepo, которые будут запущены конкуретно.
type Configuration struct {
	EventChannelSize uint64
	EventRepo        eventRepo
	EventSender      eventSender

	MaxConsumers      uint64
	ConsumerTimeout   time.Duration
	ConsumerBatchTime time.Duration
	ConsumerBatchSize uint64

	MaxProducers       uint64
	ProducerTimeout    time.Duration
	ProducerMaxWorkers uint64
}

type retranslator struct {
	consumerPool consumerPool
	producerPool producerPool
	onceStart    *sync.Once
	onceStop     *sync.Once
}

// NewRetranslator создает новый ретранслятор.
func NewRetranslator(cfg *Configuration) *retranslator {
	eventsChannel := make(chan []subscription.ServiceEvent, cfg.EventChannelSize)
	channelLocator := channellocator.NewChannelLocator(eventsChannel)

	consumerPool := consumerpool.NewConsumerPool(
		cfg.MaxConsumers,
		consumer.NewConsumerFactory(cfg.ConsumerBatchTime, cfg.ConsumerBatchSize, channelLocator, cfg.EventRepo),
		cfg.ConsumerTimeout,
	)

	producerPool := producerpool.NewProducerPool(
		cfg.MaxProducers,
		producer.NewProducerFactory(channelLocator, cfg.EventSender, cfg.EventRepo, cfg.ProducerMaxWorkers),
		cfg.ProducerTimeout,
	)

	return &retranslator{
		consumerPool: consumerPool,
		producerPool: producerPool,
		onceStart:    &sync.Once{},
		onceStop:     &sync.Once{},
	}
}

// Start запускает работу ретранслятора.
func (r *retranslator) Start(ctx context.Context) {
	r.onceStart.Do(func() {
		doneChannel := r.consumerPool.Start(ctx)
		go func() {
			ctxP, cancelCtx := context.WithCancel(context.Background())
			r.producerPool.Start(ctxP)

			<-doneChannel
			cancelCtx()
		}()
	})
}

// StopWait отправляет команду Stop пулам консьюмеров и продьюсеров,
// дожидается окончания их работы и останавливает работу ретранслятора.
//
// Обратите внимание! Метод возвращает return после остановки ретранслятора.
func (r *retranslator) StopWait() {
	r.onceStop.Do(func() {
		r.consumerPool.StopWait()
		r.producerPool.StopWait()
	})
}
