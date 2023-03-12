package domain

import (
	"route256/libs/lomsClient"
	"route256/libs/productServiceClient"
)

type CartItemAdder interface {
	AddToCart(user int64, sku uint32, count uint16) error
}

type CartListItemGetter interface {
	GetListItems(user int64) (Cart, error)
}

type CartItemDeleting interface {
	DeleteFromCart(user int64, sku uint32, count uint16) error
}

type PurchaseMaker interface {
	Purchase(user int64) error
}

type Domain struct {
	lomsClient           lomsClient.Client
	productServiceClient productServiceClient.Client
}

func New(
	lomsClient lomsClient.Client,
	productServiceClient productServiceClient.Client,
) *Domain {
	return &Domain{
		lomsClient,
		productServiceClient,
	}
}
