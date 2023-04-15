package receivers

import (
	"encoding/json"
	"fmt"
)

import (
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"route256/libs/logger"
)

var orderStatusHandler ConsumerHandlerFunc = func(msg *sarama.ConsumerMessage) error {
	var orderStatusNotification OrderStatusNotification
	err := json.Unmarshal(msg.Value, &orderStatusNotification)
	if err != nil {
		return err
	}

	logger.Info("Notification of order # "+
		fmt.Sprint(orderStatusNotification.OrderID)+
		": Order status has been changed to "+
		orderStatusNotification.Status,
		zap.Any("msg", msg),
	)

	return nil
}
