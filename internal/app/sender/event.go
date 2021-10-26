package sender

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// EventSender - интерфейс сендера.
//
// Send: отправляет событие.
type EventSender interface {
	Send(ctx context.Context, serviceEvent *subscription.ServiceEvent) error
}
