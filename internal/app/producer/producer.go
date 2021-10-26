package producer

import (
	"context"
	"log"
	"sync"
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
	Unlock(eventIDs []uint64) error
	Remove(eventIDs []uint64) error
}

type producer struct {
	timeout       time.Duration
	eventsChannel <-chan subscription.ServiceEvent
	sender        eventSender
	eventRepo     eventRepoUnlockRemover
	doneChannel   chan interface{}
	stopChannel   chan interface{}
	workerPool    *workerpool.WorkerPool
	maxWorkers    int
	onceStart     *sync.Once
	onceStop      *sync.Once
}

// NewProducer создает нового воркера-продьюсера.
//
// timeout: определяет максимальное время "пустого простоя"
// экземпляра воркера-продьюсера, по истечении которого он будет остановлен.
//
// eventsChannel: канал для чтения событий продьюсером.
//
// sender: указатель на экземпляр сендера, куда продьюсер перенаправляет
// события из канала eventsChannel.
//
// eventRepo: указатель на экземпляр репозитория событий, в котором продьюсер
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

	if maxWorkers == 0 {
		log.Panicln("maxWorkers must be greater than 0")
	}

	return &producer{
		timeout:       timeout,
		eventsChannel: eventsChannel,
		sender:        sender,
		eventRepo:     eventRepo,
		maxWorkers:    int(maxWorkers),
		onceStart:     &sync.Once{},
		onceStop:      &sync.Once{},
	}
}

// Start запускает работу продьюсера.
//
// Возвращает канал для чтения doneChannel, который закрывается продьюсером при его остановке.
func (p *producer) Start(ctx context.Context) (doneChannel <-chan interface{}) {
	p.onceStart.Do(func() {
		p.doneChannel = make(chan interface{})
		p.stopChannel = make(chan interface{})
		p.workerPool = workerpool.New(p.maxWorkers)

		sendEventAndUnlockOrRemove := func(event *subscription.ServiceEvent) {
			if err := p.sender.Send(ctx, event); err != nil {
				log.Printf("producer: failed to send event with ID %v - %v", event.ID, err)
				p.workerPool.Submit(func() {
					if err := p.eventRepo.Unlock([]uint64{event.ID}); err != nil {
						log.Printf("producer: failed to unlock event with ID %v after fail send - %v", event.ID, err)
					}
				})

				return
			}

			p.workerPool.Submit(func() {
				if err := p.eventRepo.Remove([]uint64{event.ID}); err != nil {
					log.Printf("producer: failed to remove event with ID %v after send - %v", event.ID, err)
				}
			})
		}

		go func() {
			defer func() {
				p.workerPool.StopWait()
				close(p.doneChannel)
			}()

			timeout := time.NewTimer(p.timeout)

			for {
				select {
				case event := <-p.eventsChannel:
					sendEventAndUnlockOrRemove(&event)

					continue
				case <-ctx.Done():
					p.stop()
				case <-p.stopChannel:
					if len(p.eventsChannel) == 0 {
						return
					}
				case <-timeout.C:
					return
				}
			}
		}()
	})

	return p.doneChannel
}

func (p *producer) stop() {
	p.onceStop.Do(func() {
		close(p.stopChannel)
	})
}

// StopWait отправляет команду Stop продьюсеру,
// дожидается окончания его работы и останавливает работу продьюсера.
//
// Обратите внимание! Метод возвращает return после остановки продьюсера.
func (p *producer) StopWait() {
	p.stop()
	<-p.doneChannel
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
// sender: указатель на экземпляр сендера, куда продьюсер перенаправляет
// события из канала eventsChannel.
//
// eventRepo: указатель на экземпляр репозитория событий, в котором продьюсер
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
