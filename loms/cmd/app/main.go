package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"route256/loms/internal/client/kafka"
	"route256/loms/internal/producers"
)

import (
	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"route256/libs/config"
	"route256/libs/db"
	"route256/libs/db/transaction"
	"route256/loms/internal/api/lomsV1"
	AppConfig "route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/repositories"
	"route256/loms/internal/scheduler/cron/jobs"
	desc "route256/loms/pkg/loms_v1"
)

import (
	_ "github.com/lib/pq"
)

type Logger struct{}

func (l *Logger) Printf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	fmt.Println()
}

func main() {
	configApp := &AppConfig.Config{}
	err := config.Init("config.yml", configApp)
	if err != nil {
		log.Fatal("config init", err)
	}

	ctx := context.Background()

	// Создание клиента для работы с DB
	clientDb, err := db.NewClient(ctx, &db.Config{
		Host:     configApp.DataBase.Host,
		Port:     configApp.DataBase.Port,
		User:     configApp.DataBase.User,
		Password: configApp.DataBase.Password,
		Name:     configApp.DataBase.Name,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer clientDb.Close()

	// TransactionManager
	transactionManager := transaction.NewTransactionManager(clientDb.DB())

	// Repositories
	stockRepository := repositories.NewStockRepository(
		transactionManager,
	)
	orderRepository := repositories.NewOrderRepository(
		transactionManager,
	)
	orderItemRepository := repositories.NewOrderItemRepository(
		transactionManager,
	)
	orderItemStockRepository := repositories.NewOrderItemStockRepository(
		transactionManager,
	)

	// Producers
	producer, err := kafka.NewSyncProducer(configApp.Brokers)
	if err != nil {
		log.Fatal(err)
	}
	orderStatusNotifier := producers.NewOrderStatusNotifier(producer, "order_statuses")

	// Create domain
	domain := domain.NewDomain(
		transactionManager,
		stockRepository,
		orderRepository,
		orderItemRepository,
		orderItemStockRepository,
		orderStatusNotifier,
	)

	// CRON
	logger := cron.VerbosePrintfLogger(&Logger{})
	c := cron.New(
		cron.WithLogger(logger),
	)
	defer c.Stop()

	c.AddJob("* * * * *", jobs.NewOrdersChecker(domain, logger))
	c.Start()

	// Create tcp listener
	lis, err := net.Listen("tcp", ":"+configApp.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	reflection.Register(server)
	desc.RegisterLomsV1Server(server, lomsV1.NewLomsV1(domain))

	log.Printf("server listening at %v", lis.Addr())

	// Run server
	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
