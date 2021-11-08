package channellocator

import (
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// ChannelLocator предоставляет интерфейс-локатор каналов обмена в ретрансляторе.
type ChannelLocator interface {
	GetMainEventsReadChannel() <-chan []subscription.ServiceEvent
	GetMainEventsWriteChannel() chan<- []subscription.ServiceEvent
}

type channelLocator struct {
	mainEventsChannel chan []subscription.ServiceEvent
}

// NewChannelLocator создает новый локатор каналов.
//
// mainEventsChannel - основной канал обмена событиями.
func NewChannelLocator(mainEventsChannel chan []subscription.ServiceEvent) *channelLocator {
	return &channelLocator{
		mainEventsChannel: mainEventsChannel,
	}
}

func (cl *channelLocator) GetMainEventsReadChannel() <-chan []subscription.ServiceEvent {
	return cl.mainEventsChannel
}

func (cl *channelLocator) GetMainEventsWriteChannel() chan<- []subscription.ServiceEvent {
	return cl.mainEventsChannel
}
