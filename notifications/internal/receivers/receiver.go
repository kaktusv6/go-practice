package receivers

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

type Receiver struct {
	consumer sarama.Consumer
	handlers Handlers
}

func NewReceiver(
	consumer sarama.Consumer,
	handlers Handlers,
) Receiver {
	return Receiver{
		consumer: consumer,
		handlers: handlers,
	}
}

func (r *Receiver) Subscribe(topic string) error {
	handler, ok := r.handlers[topic]
	if !ok {
		return errors.New("No exist handler for topic: " + topic)
	}

	partitions, err := r.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := r.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				handler(message)
			}
		}(pc)
	}

	return nil
}
