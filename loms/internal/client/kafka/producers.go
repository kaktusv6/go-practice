package kafka

import "github.com/Shopify/sarama"

func NewSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
