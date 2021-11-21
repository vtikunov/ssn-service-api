package kafka

import (
	"context"

	"github.com/Shopify/sarama"

	"github.com/ozonmp/ssn-service-api/internal/pkg/logger"
)

// ConsumeFunction - функция для передачи стартеру.
type ConsumeFunction func(ctx context.Context, message *sarama.ConsumerMessage) error

type consumer struct {
	fn ConsumeFunction
}

func (consumer *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := consumer.fn(session.Context(), message)

		if err != nil {
			return err
		}

		session.MarkMessage(message, "")
	}

	return nil
}

// StartConsuming - функция-стартер консьюмера.
func StartConsuming(ctx context.Context, brokers []string, topic string, group string, consumeFunction ConsumeFunction) error {

	config := sarama.NewConfig()

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)

	if err != nil {
		return err
	}

	cns := consumer{
		fn: consumeFunction,
	}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, &cns); err != nil {
				logger.ErrorKV(ctx, "error from consumer", "err", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}
