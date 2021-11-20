package sender

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/ozonmp/ssn-service-api/internal/model/subscription"
	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"

	pb "github.com/ozonmp/ssn-service-api/pkg/ssn-service-api"
)

// EventSender - интерфейс сендера.
//
// Send: отправляет событие.
type EventSender interface {
	Send(event *subscription.ServiceEvent) error
}

type kafkaSender struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaSender - создает сендера событий в Apache Kafka.
func NewKafkaSender(ctx context.Context, brokers []string, topic string, partitionFactor uint8) *kafkaSender {
	config := sarama.NewConfig()
	config.Producer.Partitioner = newPartitionerConstructor(partitionFactor)
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		logger.FatalKV(ctx, "failed init kafka producer", "err", err)
	}

	return &kafkaSender{
		producer: producer,
		topic:    topic,
	}
}

func (s *kafkaSender) Send(event *subscription.ServiceEvent) error {
	pbSrvPayload := &pb.ServiceEventPayload{}

	if event.Service != nil {
		pbSrvPayload.ServiceId = event.Service.ID
		pbSrvPayload.Name = event.Service.Name
		pbSrvPayload.Description = event.Service.Description
	}

	pbEvent := &pb.ServiceEvent{
		Id:        event.ID,
		ServiceId: event.ServiceID,
		Type:      string(event.Type),
		Subtype:   string(event.SubType),
		Payload:   pbSrvPayload,
	}

	msgValue, err := proto.Marshal(pbEvent)

	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Key:   &key{serviceID: event.ServiceID},
		Topic: s.topic,
		Value: sarama.ByteEncoder(msgValue),
	})

	return err
}

type partitioner struct {
	partitionFactor uint8
}

func newPartitionerConstructor(partitionFactor uint8) sarama.PartitionerConstructor {
	return func(topic string) sarama.Partitioner {
		return &partitioner{
			partitionFactor: partitionFactor,
		}
	}
}

func (p *partitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	if int32(p.partitionFactor) != numPartitions {
		return 0, fmt.Errorf("numPartitions %v must be equal to partitionFactor %v", numPartitions, p.partitionFactor)
	}

	b, err := message.Key.Encode()

	if err != nil {
		return 0, err
	}

	serviceID := binary.LittleEndian.Uint64(b)

	return int32(serviceID % uint64(p.partitionFactor)), nil
}

func (p *partitioner) RequiresConsistency() bool {
	return true
}

type key struct {
	serviceID uint64
}

func (k *key) Encode() ([]byte, error) {
	return k.encode(), nil
}

func (k *key) Length() int {
	return len(k.encode())
}

func (k *key) encode() []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k.serviceID)

	return b
}
