package main

import (
	"context"
	"log"
	"route256/libs/config"
	"route256/notifications/internal/clients/kafka"
	AppConfig "route256/notifications/internal/config"
	"route256/notifications/internal/receivers"
)

func main() {
	configApp := &AppConfig.Config{}
	err := config.Init("config.yml", configApp)
	if err != nil {
		log.Fatal("config init", err)
	}

	consumer, err := kafka.NewConsumer(configApp.Brokers)
	if err != nil {
		log.Fatalln(err)
	}

	handlers := receivers.InitHandlers()

	receiver := receivers.NewReceiver(consumer, handlers)
	err = receiver.Subscribe("order_statuses")
	if err != nil {
		log.Fatal("Receiver error", err)
	}

	<-context.TODO().Done()
}
