package producers

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"route256/loms/internal/domain"
	"time"
)

type orderStatusNotifier struct {
	producer sarama.SyncProducer
	topic    string
}

func NewOrderStatusNotifier(
	producer sarama.SyncProducer,
	topic string,
) domain.OrderStatusNotifier {
	return &orderStatusNotifier{
		producer: producer,
		topic:    topic,
	}
}

func (o *orderStatusNotifier) Notify(order *domain.Order) error {
	jsonBytes, err := json.Marshal(OrderStatusNotification{
		OrderID: order.ID,
		Status:  order.Status,
	})
	if err != nil {
		return err
	}

	_, _, err = o.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     o.topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(jsonBytes),
		Key:       sarama.StringEncoder(fmt.Sprint(order.ID)),
		Timestamp: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
