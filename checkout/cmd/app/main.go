package main

import (
	"context"
	"log"
	"net"
	"net/http"
)

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	checkoutV1 "route256/checkout/internal/api/checkout_v1"
	AppConfig "route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repositories"
	desc "route256/checkout/pkg/checkout_v1"
	productServiceV1Clinet "route256/checkout/pkg/product_service_v1"
	"route256/libs/config"
	"route256/libs/logger"
	"route256/libs/metrics"
	"route256/libs/tracing"

	"route256/libs/db"
	"route256/libs/db/transaction"
	lomsV1Clinet "route256/loms/pkg/loms_v1"
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

	lis, err := net.Listen("tcp", ":"+configApp.App.Port)
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
		panic(err)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			metrics.Metrics,
			tracing.Tracer,
		),
	)

	lomsCon, err := grpc.Dial(
		configApp.Loms.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		logger.Fatal("Error create connection to loms", zap.Error(err))
		panic(err)
	}
	defer lomsCon.Close()

	productServiceCon, err := grpc.Dial(
		configApp.ProductService.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		logger.Fatal("Error create connection to productService", zap.Error(err))
		panic(err)
	}
	defer productServiceCon.Close()

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

	provider := transaction.NewTransactionManager(clientDb.DB())

	cartItemRepository := repositories.NewOrderItemRepository(provider)

	productRepository := repositories.NewOrderProductRepository(
		productServiceV1Clinet.NewProductServiceClient(productServiceCon),
		configApp.ProductService.Token,
	)

	stockRepository := repositories.NewStockRepository(
		lomsV1Clinet.NewLomsV1Client(lomsCon),
	)
	orderRepository := repositories.NewOrderRepository(
		lomsV1Clinet.NewLomsV1Client(lomsCon),
	)

	domain := domain.New(
		stockRepository,
		orderRepository,
		provider,
		cartItemRepository,
		productRepository,
	)

	reflection.Register(server)
	desc.RegisterCheckoutV1Server(server, checkoutV1.New(domain))

	logger.Info("server listening at", zap.Any("address", lis.Addr()))

	grpc_prometheus.Register(server)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		errMetrics := http.ListenAndServe(":8080", nil)
		logger.Fatal("Error list metrics", zap.Error(errMetrics))
		panic(errMetrics)
	}()

	if err = server.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
		panic(err)
	}
}
