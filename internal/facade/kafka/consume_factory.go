package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
)

type serviceService interface {
	Add(ctx context.Context, service *subscription.Service) error
	Update(ctx context.Context, service *subscription.Service) error
	Remove(ctx context.Context, serviceID uint64, eventID uint64) error
}

// GetServiceEventConsume - фабрика функций ConsumeFunction.
func GetServiceEventConsume(srvService serviceService) ConsumeFunction {
	return func(ctx context.Context, message *sarama.ConsumerMessage) error {
		var event pb.ServiceEvent

		err := proto.Unmarshal(message.Value, &event)

		if err != nil {
			return err
		}

		serviceEvent := convertPbToEvent(&event)

		switch serviceEvent.Type {
		case subscription.Created:
			return srvService.Add(ctx, serviceEvent.Service)
		case subscription.Updated:
			return srvService.Update(ctx, serviceEvent.Service)
		case subscription.Removed:
			return srvService.Remove(ctx, serviceEvent.ServiceID, serviceEvent.ID)
		}

		return fmt.Errorf("unknown event type \"%v\" at event ID %v", serviceEvent.Type, serviceEvent.ID)
	}
}

func convertPbToEvent(pb *pb.ServiceEvent) *subscription.ServiceEvent {
	var service *subscription.Service

	if pb.Payload != nil {
		service = &subscription.Service{
			ID:          pb.Payload.GetServiceId(),
			Name:        pb.Payload.GetName(),
			Description: pb.Payload.GetDescription(),
			LastEventID: pb.GetId(),
		}
	}

	return &subscription.ServiceEvent{
		ID:        pb.GetId(),
		ServiceID: pb.GetServiceId(),
		Type:      subscription.EventType(pb.GetType()),
		SubType:   subscription.EventSubType(pb.GetSubtype()),
		Service:   service,
	}
}
