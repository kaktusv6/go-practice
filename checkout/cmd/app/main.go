package main

import (
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
	desc "route256/checkout/pkg/checkout_v1"
	"route256/libs/config"
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

	domain := domain.New(lomsCon, productServiceCon, configApp.ProductService.Token)

	reflection.Register(server)
	desc.RegisterCheckoutV1Server(server, checkoutV1.New(domain))

	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
