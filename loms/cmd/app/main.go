package main

import (
	"log"
	"net/http"
)

import (
	"route256/libs/config"
	"route256/libs/httpServerWrapper"
	AppConfig "route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/handlers/cancelOrder"
	"route256/loms/internal/handlers/createOrder"
	"route256/loms/internal/handlers/listOrder"
	"route256/loms/internal/handlers/orderPayed"
	"route256/loms/internal/handlers/stocks"
)

func main() {
	configApp := &AppConfig.Config{}
	err := config.Init("config.yml", configApp)
	if err != nil {
		log.Fatal("config init", err)
	}

	domain := domain.NewDomain()

	stockHandler := stocks.New(domain)
	cancelOrderHandler := cancelOrder.New(domain)
	orderPayedHandler := orderPayed.New(domain)
	listOrderHandler := listOrder.New(domain)
	createOrderHandler := createOrder.New(domain)

	http.Handle("/stocks", httpServerWrapper.New(stockHandler.Handle))
	http.Handle("/cancelOrder", httpServerWrapper.New(cancelOrderHandler.Handle))
	http.Handle("/orderPayed", httpServerWrapper.New(orderPayedHandler.Handle))
	http.Handle("/listOrder", httpServerWrapper.New(listOrderHandler.Handle))
	http.Handle("/createOrder", httpServerWrapper.New(createOrderHandler.Handle))

	log.Println("Listening HTTP at", configApp.App.Port)
	err = http.ListenAndServe(":"+configApp.App.Port, nil)
	log.Fatal("Error listen HTTP", err)
}
