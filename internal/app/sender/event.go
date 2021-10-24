package sender

import (
	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// EventSender - интерфейс сендера.
//
// Send: отправляет событие.
type EventSender interface {
	Send(serviceEvent *subscription.ServiceEvent) error
}
