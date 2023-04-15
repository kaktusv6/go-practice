package main

import (
	"context"
	"log"
	"net"
	"net/http"
)

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"route256/libs/config"
	"route256/libs/db"
	"route256/libs/db/transaction"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracing"
	"route256/loms/internal/api/lomsV1"
	"route256/loms/internal/client/kafka"
	AppConfig "route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/producers"
	"route256/loms/internal/repositories"
	"route256/loms/internal/scheduler/cron/jobs"
	desc "route256/loms/pkg/loms_v1"
)

import (
	_ "github.com/lib/pq"
)

func main() {
	configApp := &AppConfig.Config{}
	err := config.Init("config.yml", configApp)
	if err != nil {
		log.Fatal("config init", err)
	}

	loggerConfig := logger.Config{
		Level: configApp.Logger.Level,
		Env:   configApp.App.Environment,
	}
	logger.Init(loggerConfig)

	metrics.Init("hw")

	tracing.Init(configApp.App.Name)

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
		logger.Fatal(err.Error())
		panic(err)
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
		logger.Fatal(err.Error())
		panic(err)
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
	c := cron.New()
	defer c.Stop()

	c.AddJob("* * * * *", jobs.NewOrdersChecker(domain))
	c.Start()

	// Create tcp listener
	lis, err := net.Listen("tcp", ":"+configApp.App.Port)
	if err != nil {
		logger.Fatal(err.Error())
		panic(err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			metrics.Metrics,
			tracing.Tracer,
		),
	)

	reflection.Register(server)

	desc.RegisterLomsV1Server(server, lomsV1.NewLomsV1(domain))

	logger.Info("server listening at", zap.String("address", lis.Addr().String()))

	grpc_prometheus.Register(server)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		errMetrics := http.ListenAndServe(":8080", nil)
		logger.Fatal("Error list metrics", zap.Error(errMetrics))
		panic(errMetrics)
	}()

	// Run server
	if err = server.Serve(lis); err != nil {
		logger.Fatal("failed to serve: ", zap.Error(err))
		panic(err)
	}
}
