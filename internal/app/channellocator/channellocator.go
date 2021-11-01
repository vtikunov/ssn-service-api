package channellocator

import (
	"errors"
	"sync"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// ChannelLocator предоставляет интерфейс-локатор каналов обмена в ретрансляторе.
type ChannelLocator interface {
	GetMainEventsReadChannel() <-chan []subscription.ServiceEvent
	GetMainEventsWriteChannel() chan<- []subscription.ServiceEvent
	GetEventsServiceIDWriteChannel(serviceID uint64) chan<- []subscription.ServiceEvent
	GetEventsServiceIDReadChannel(serviceID uint64) (<-chan []subscription.ServiceEvent, error)
}

type channelLocator struct {
	mainEventsChannel chan []subscription.ServiceEvent
	channelsMap       *sync.Map
}

// NewChannelLocator создает новый локатор каналов.
//
// mainEventsChannel - основной канал обмена событиями.
func NewChannelLocator(mainEventsChannel chan []subscription.ServiceEvent) *channelLocator {
	return &channelLocator{
		mainEventsChannel: mainEventsChannel,
		channelsMap:       &sync.Map{},
	}
}

func (cl *channelLocator) GetMainEventsReadChannel() <-chan []subscription.ServiceEvent {
	return cl.mainEventsChannel
}

func (cl *channelLocator) GetMainEventsWriteChannel() chan<- []subscription.ServiceEvent {
	return cl.mainEventsChannel
}

func (cl *channelLocator) GetEventsServiceIDWriteChannel(serviceID uint64) chan<- []subscription.ServiceEvent {
	channel, _ := cl.channelsMap.LoadOrStore(serviceID, make(chan []subscription.ServiceEvent))

	return channel.(chan []subscription.ServiceEvent)
}

func (cl *channelLocator) GetEventsServiceIDReadChannel(serviceID uint64) (<-chan []subscription.ServiceEvent, error) {
	channel, ok := cl.channelsMap.Load(serviceID)
	if !ok {
		return nil, errors.New("channel not found in locator")
	}

	return channel.(chan []subscription.ServiceEvent), nil
}
