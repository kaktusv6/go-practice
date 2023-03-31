package main

import (
	"context"
	"log"
	"net"
)

import (
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

	lis, err := net.Listen("tcp", ":"+configApp.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	lomsCon, err := grpc.Dial(configApp.Loms.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error create connection to loms", err)
	}
	defer lomsCon.Close()

	productServiceCon, err := grpc.Dial(configApp.ProductService.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error create connection to productService", err)
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

	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
