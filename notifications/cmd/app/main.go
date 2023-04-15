package main

import (
	"context"
	"log"
)

import (
	"go.uber.org/zap"
	"route256/libs/config"
	"route256/libs/logger"
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
	configLogger := logger.Config{
		Env:   configApp.App.Env,
		Level: configApp.Logger.Level,
	}
	logger.Init(configLogger)

	consumer, err := kafka.NewConsumer(configApp.Brokers)
	if err != nil {
		logger.Fatal("Error create consumer", zap.Error(err))
	}

	handlers := receivers.InitHandlers()

	receiver := receivers.NewReceiver(consumer, handlers)
	err = receiver.Subscribe("order_statuses")
	if err != nil {
		logger.Fatal("Receiver error", zap.Error(err))
	}

	<-context.TODO().Done()
}
