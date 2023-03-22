package domain

import (
	"context"
)

import (
	"google.golang.org/grpc"
	"route256/libs/transactor"
	lomsV1Clinet "route256/loms/pkg/loms_v1"
)

type Domain interface {
	AddToCart(ctx context.Context, itemInfo CartItem) error
	GetListItems(ctx context.Context, user int64) (*Cart, error)
	DeleteFromCart(ctx context.Context, itemInfo CartItem) error
	Purchase(ctx context.Context, user int64) error
}

type domain struct {
	lomsClient         lomsV1Clinet.LomsV1Client
	transactionManager *transactor.TransactionManager
	cartItemRepository CartItemRepository
	productRepository  ProductRepository
}

func New(
	lomsConnection *grpc.ClientConn,
	transactionManager *transactor.TransactionManager,
	cartItemRepository CartItemRepository,
	productRepository ProductRepository,
) Domain {
	return &domain{
		lomsV1Clinet.NewLomsV1Client(lomsConnection),
		transactionManager,
		cartItemRepository,
		productRepository,
	}
}
