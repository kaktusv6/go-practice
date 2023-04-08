package receivers

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
)

var orderStatusHandler ConsumerHandlerFunc = func(msg *sarama.ConsumerMessage) error {
	var orderStatusNotification OrderStatusNotification
	err := json.Unmarshal(msg.Value, &orderStatusNotification)
	if err != nil {
		return err
	}

	log.Printf("Notification of order # %d: Order status has been changed to %s", orderStatusNotification.OrderID, orderStatusNotification.Status)

	return nil
}
