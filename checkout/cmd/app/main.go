package main

import (
	"log"
	"net/http"
)

import (
	AppConfig "route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addToCart"
	"route256/checkout/internal/handlers/deleteFromCart"
	"route256/checkout/internal/handlers/listCart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/config"
	"route256/libs/httpServerWrapper"
	"route256/libs/lomsClient"
	"route256/libs/productServiceClient"
)

func main() {
	configApp := &AppConfig.Config{}
	err := config.Init("config.yml", configApp)
	if err != nil {
		log.Fatal("config init", err)
	}

	lomsClient := lomsClient.New(
		lomsClient.NewConfig(
			configApp.Loms.Url,
		))

	prodServClient := productServiceClient.New(
		productServiceClient.NewConfig(
			configApp.ProductService.Url,
			configApp.ProductService.Token,
		))

	domain := domain.New(lomsClient, prodServClient)

	addToCartHandler := addToCart.New(domain)
	listCartHandler := listCart.New(domain)
	purchaseHandler := purchase.New(domain)
	deleteFromCartHandler := deleteFromCart.New(domain)

	http.Handle("/addToCart", httpServerWrapper.New(addToCartHandler.Handle))
	http.Handle("/listCart", httpServerWrapper.New(listCartHandler.Handle))
	http.Handle("/purchase", httpServerWrapper.New(purchaseHandler.Handle))
	http.Handle("/deleteFromCart", httpServerWrapper.New(deleteFromCartHandler.Handle))

	log.Println("Listening HTTP at", configApp.App.Port)
	err = http.ListenAndServe(":"+configApp.App.Port, nil)
	log.Fatal("Error listen HTTP", err)
}
