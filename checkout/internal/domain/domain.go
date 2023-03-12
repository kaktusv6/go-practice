package domain

import (
	"context"
)

import (
	"google.golang.org/grpc"
	productServiceV1Clinet "route256/checkout/pkg/product_service_v1"
	lomsV1Clinet "route256/loms/pkg/loms_v1"
)

type Domain interface {
	AddToCart(ctx context.Context, itemInfo ItemInfo) error
	GetListItems(ctx context.Context, user int64) (Cart, error)
	DeleteFromCart(ctx context.Context, itemInfo ItemInfo) error
	Purchase(ctx context.Context, user int64) error
}

type domain struct {
	lomsClient           lomsV1Clinet.LomsV1Client
	productServiceClient productServiceV1Clinet.ProductServiceClient
	productServiceToken  string
}

func New(
	lomsConnection *grpc.ClientConn,
	productServiceConnection *grpc.ClientConn,
	productServiceToken string,
) Domain {
	return &domain{
		lomsV1Clinet.NewLomsV1Client(lomsConnection),
		productServiceV1Clinet.NewProductServiceClient(productServiceConnection),
		productServiceToken,
	}
}
