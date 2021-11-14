package retranslator

import (
	"context"
	"sync"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/channellocator"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/config"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/consumer"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/consumerpool"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/producer"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/producerpool"
	"github.com/ozonmp/ssn-service-api/internal/retranslator/repo"
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
	Lock(ctx context.Context, n uint64, tx repo.QueryerExecer) ([]subscription.ServiceEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64, tx repo.QueryerExecer) error

	Remove(ctx context.Context, eventIDs []uint64, tx repo.QueryerExecer) error
}

type eventSender interface {
	Send(ctx context.Context, serviceEvent *subscription.ServiceEvent) error
}

type retranslator struct {
	consumerPool consumerPool
	producerPool producerPool
	onceStart    *sync.Once
	onceStop     *sync.Once
}

// NewRetranslator создает новый ретранслятор.
func NewRetranslator(ctx context.Context, cfg *config.Retranslator, repo eventRepo, sender eventSender) *retranslator {
	eventsChannel := make(chan []subscription.ServiceEvent, cfg.EventChannelSize)
	channelLocator := channellocator.NewChannelLocator(eventsChannel)

	consumerPool := consumerpool.NewConsumerPool(
		ctx,
		cfg.MaxConsumers,
		consumer.NewConsumerFactory(cfg.ConsumerBatchTime, cfg.ConsumerBatchSize, channelLocator, repo),
		cfg.ConsumerTimeout,
	)

	producerPool := producerpool.NewProducerPool(
		ctx,
		cfg.MaxProducers,
		producer.NewProducerFactory(channelLocator, sender, repo, cfg.ProducerMaxWorkers),
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
