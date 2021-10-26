package producer

import (
	"context"
	"log"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// Producer - общий интерфейс воркеров-продьюсеров для совместимости между пакетами.
type Producer interface {
	Start(ctx context.Context) (doneChannel <-chan interface{})
	StopWait()
}

type eventSender interface {
	Send(ctx context.Context, serviceEvent *subscription.ServiceEvent) error
}

type eventRepoUnlockRemover interface {
	Unlock(ctx context.Context, eventIDs []uint64) error
	Remove(ctx context.Context, eventIDs []uint64) error
}

type producer struct {
	timeout       time.Duration
	eventsChannel <-chan subscription.ServiceEvent
	sender        eventSender
	eventRepo     eventRepoUnlockRemover
	doneChannel   chan interface{}
	stopChannel   chan interface{}
	isStarted     bool
	workerPool    *workerpool.WorkerPool
	maxWorkers    int
}

// NewProducer создает нового воркера-продьюсера.
//
// timeout: определяет максимальное время "пустого простоя"
// экземпляра воркера-продьюсера, по истечении которого он будет остановлен.
//
// eventsChannel: канал для чтения событий продьюсером.
//
// sender: ссылка на экземпляр сендера, куда продьюсер перенаправляет
// события из канала eventsChannel.
//
// eventRepo: ссылка на экземпляр репозитория событий, в котором продьюсер
// разблокирует или удаляет события после перенаправления.
//
// maxWorkers: определяет максимальное количество вспомогательных воркеров работы
// с репозиторием событий eventRepo, которые будут запущены конкуретно.
func NewProducer(
	timeout time.Duration,
	eventsChannel <-chan subscription.ServiceEvent,
	sender eventSender,
	eventRepo eventRepoUnlockRemover,
	maxWorkers uint64,
) *producer {

	return &producer{
		timeout:       timeout,
		eventsChannel: eventsChannel,
		sender:        sender,
		eventRepo:     eventRepo,
		maxWorkers:    int(maxWorkers),
	}
}

// Start запускает работу продьюсера.
//
// Возвращает канал для чтения doneChannel, который закрывается продьюсером при его остановке.
func (p *producer) Start(ctx context.Context) (doneChannel <-chan interface{}) {
	if p.isStarted {
		log.Panic("producer is already started")
	}
	p.isStarted = true
	p.doneChannel = make(chan interface{})
	p.stopChannel = make(chan interface{})
	p.workerPool = workerpool.New(p.maxWorkers)

	sendEventAndUnlockOrRemove := func(event *subscription.ServiceEvent) {
		if err := p.sender.Send(ctx, event); err != nil {
			log.Printf("producer: failed to send event with ID %v - %v", event.ID, err)
			p.workerPool.Submit(func() {
				if err := p.eventRepo.Unlock(ctx, []uint64{event.ID}); err != nil {
					log.Printf("producer: failed to unlock event with ID %v after fail send - %v", event.ID, err)
				}
			})

			return
		}

		p.workerPool.Submit(func() {
			if err := p.eventRepo.Remove(ctx, []uint64{event.ID}); err != nil {
				log.Printf("producer: failed to remove event with ID %v after send - %v", event.ID, err)
			}
		})
	}

	go func() {
		defer close(p.doneChannel)
		defer p.workerPool.StopWait()
		timeout := time.NewTimer(p.timeout)

		for {
			select {
			case event := <-p.eventsChannel:
				sendEventAndUnlockOrRemove(&event)
			case <-ctx.Done():
				if len(p.eventsChannel) == 0 {
					return
				}
			case <-p.stopChannel:
				if len(p.eventsChannel) == 0 {
					return
				}
			case <-timeout.C:
				return
			}
		}
	}()

	return p.doneChannel
}

// StopWait отправляет команду Stop продьюсеру,
// дожидается окончания его работы и останавливает работу продьюсера.
//
// Обратите внимание! Метод возвращает return после остановки продьюсера.
func (p *producer) StopWait() {
	if !p.isStarted {
		log.Panic("producer is already stopped")
	}
	close(p.stopChannel)
	<-p.doneChannel

	p.isStarted = false
}

type producerFactory struct {
	eventsChannel <-chan subscription.ServiceEvent
	sender        eventSender
	eventRepo     eventRepoUnlockRemover
	maxWorkers    uint64
}

// NewProducerFactory создает фабрику воркеров-продьюсеров.
//
// eventsChannel: канал для чтения событий продьюсером.
//
// sender: ссылка на экземпляр сендера, куда продьюсер перенаправляет
// события из канала eventsChannel.
//
// eventRepo: ссылка на экземпляр репозитория событий, в котором продьюсер
// разблокирует или удаляет события после перенаправления.
//
// maxWorkers: определяет максимальное количество вспомогательных воркеров работы
// с репозиторием событий eventRepo, которые будут запущены конкуретно.
func NewProducerFactory(
	eventsChannel <-chan subscription.ServiceEvent,
	sender eventSender,
	eventRepo eventRepoUnlockRemover,
	maxWorkers uint64,
) *producerFactory {

	return &producerFactory{
		eventsChannel: eventsChannel,
		sender:        sender,
		eventRepo:     eventRepo,
		maxWorkers:    maxWorkers,
	}
}

// Create создает воркера-продьюсера
//
// timeout: определяет максимальное время "пустого простоя" экземпляра воркера-продьюсера,
// по истечении которого он будет остановлен.
func (pf *producerFactory) Create(timeout time.Duration) Producer {
	return NewProducer(timeout, pf.eventsChannel, pf.sender, pf.eventRepo, pf.maxWorkers)
}
