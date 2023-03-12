package main

import (
	"log"
	"net"
)

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"route256/libs/config"
	"route256/loms/internal/api/lomsV1"
	AppConfig "route256/loms/internal/config"
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
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

	domain := domain.NewDomain()

	reflection.Register(server)
	desc.RegisterLomsV1Server(server, lomsV1.NewLomsV1(domain))

	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
