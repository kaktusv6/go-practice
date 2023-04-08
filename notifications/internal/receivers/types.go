package receivers

import (
	"github.com/Shopify/sarama"
)

type ConsumerHandlerFunc func(*sarama.ConsumerMessage) error

type Handlers map[string]ConsumerHandlerFunc
