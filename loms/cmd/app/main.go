package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
)

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"route256/libs/config"
	"route256/libs/transactor"
	"route256/loms/internal/api/lomsV1"
	AppConfig "route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/repositories"
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

	// connection string
	psqlConn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configApp.DataBase.Host,
		configApp.DataBase.Port,
		configApp.DataBase.User,
		configApp.DataBase.Password,
		configApp.DataBase.Name,
	)

	ctx := context.Background()

	// open database
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	provider := transactor.NewQueryEngineProvider(pool)

	stockRepository := repositories.NewStockRepository(
		provider,
	)
	orderRepository := repositories.NewOrderRepository(
		provider,
	)
	orderItemRepository := repositories.NewOrderItemRepository(
		provider,
	)
	orderItemStockRepository := repositories.NewOrderItemStockRepository(
		provider,
	)

	domain := domain.NewDomain(
		provider,
		stockRepository,
		orderRepository,
		orderItemRepository,
		orderItemStockRepository,
	)

	lis, err := net.Listen("tcp", ":"+configApp.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	reflection.Register(server)
	desc.RegisterLomsV1Server(server, lomsV1.NewLomsV1(domain))

	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
