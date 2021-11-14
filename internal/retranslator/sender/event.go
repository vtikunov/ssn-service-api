package sender

import (
	"context"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
)

// EventSender - интерфейс сендера.
//
// Send: отправляет событие.
type EventSender interface {
	Send(ctx context.Context, serviceEvent *subscription.ServiceEvent) error
}

type dummySender struct{}

// NewDummySender - создает заглушку сендера.
func NewDummySender() *dummySender {
	return &dummySender{}
}

func (d *dummySender) Send(ctx context.Context, serviceEvent *subscription.ServiceEvent) error {
	logger.DebugKV(ctx, "event sent", "eventID", serviceEvent.ID, "serviceID", serviceEvent.ServiceID)
	return nil
}
