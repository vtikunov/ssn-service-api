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

type channelLocator interface {
	GetMainEventsReadChannel() <-chan []subscription.ServiceEvent
	GetEventsServiceIDReadChannel(serviceID uint64) (<-chan []subscription.ServiceEvent, error)
}

type producer struct {
	timeout        time.Duration
	channelLocator channelLocator
	sender         eventSender
	eventRepo      eventRepoUnlockRemover
	doneChannel    chan interface{}
	stopChannel    chan interface{}
	workerPool     *workerpool.WorkerPool
	maxWorkers     int
	onceStart      *sync.Once
	onceStop       *sync.Once
}

// NewProducer создает нового воркера-продьюсера.
//
// timeout: определяет максимальное время "пустого простоя"
// экземпляра воркера-продьюсера, по истечении которого он будет остановлен.
//
// channelLocator: локатор каналов.
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
	channelLocator channelLocator,
	sender eventSender,
	eventRepo eventRepoUnlockRemover,
	maxWorkers uint64,
) *producer {

	if maxWorkers == 0 {
		log.Panicln("maxWorkers must be greater than 0")
	}

	return &producer{
		timeout:        timeout,
		channelLocator: channelLocator,
		sender:         sender,
		eventRepo:      eventRepo,
		maxWorkers:     int(maxWorkers),
		onceStart:      &sync.Once{},
		onceStop:       &sync.Once{},
	}
}

func (p *producer) sendEventsAndUnlockOrRemove(ctx context.Context, events []subscription.ServiceEvent) {
	errorIDs := make([]uint64, 0)
	completeIDs := make([]uint64, 0)
	for _, event := range events {
		eventChannel, err := p.channelLocator.GetEventsServiceIDReadChannel(event.Service.ID)
		for err != nil {
			eventChannel, err = p.channelLocator.GetEventsServiceIDReadChannel(event.Service.ID)
		}

		eventsForServiceID := <-eventChannel

		for es, eventForService := range eventsForServiceID {
			if err := p.sender.Send(ctx, &eventsForServiceID[es]); err != nil {
				log.Printf("producer: failed to send event with ID %v - %v", eventForService.ID, err)
				errorIDs = append(errorIDs, eventForService.ID)

				continue
			}

			completeIDs = append(completeIDs, eventForService.ID)
		}
	}

	if len(errorIDs) > 0 {
		p.workerPool.Submit(func() {
			if err := p.eventRepo.Unlock(errorIDs); err != nil {
				log.Printf("producer: failed to unlock events after fail send - %v", err)
			}
		})
	}

	if len(completeIDs) > 0 {
		p.workerPool.Submit(func() {
			if err := p.eventRepo.Remove(completeIDs); err != nil {
				log.Printf("producer: failed to remove events after send - %v", err)
			}
		})
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

		go func() {
			defer func() {
				p.workerPool.StopWait()
				close(p.doneChannel)
			}()

			mainEventsChannel := p.channelLocator.GetMainEventsReadChannel()
			timeout := time.NewTimer(p.timeout)

			for {
				select {
				case events := <-mainEventsChannel:
					p.sendEventsAndUnlockOrRemove(ctx, events)

					continue
				case <-ctx.Done():
					p.stop()
				case <-p.stopChannel:
					if len(mainEventsChannel) == 0 {
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
	channelLocator channelLocator
	sender         eventSender
	eventRepo      eventRepoUnlockRemover
	maxWorkers     uint64
}

// NewProducerFactory создает фабрику воркеров-продьюсеров.
//
// channelLocator: локатор каналов.
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
	channelLocator channelLocator,
	sender eventSender,
	eventRepo eventRepoUnlockRemover,
	maxWorkers uint64,
) *producerFactory {

	return &producerFactory{
		channelLocator: channelLocator,
		sender:         sender,
		eventRepo:      eventRepo,
		maxWorkers:     maxWorkers,
	}
}

// Create создает воркера-продьюсера
//
// timeout: определяет максимальное время "пустого простоя" экземпляра воркера-продьюсера,
// по истечении которого он будет остановлен.
func (pf *producerFactory) Create(timeout time.Duration) Producer {
	return NewProducer(timeout, pf.channelLocator, pf.sender, pf.eventRepo, pf.maxWorkers)
}
