package domain

import (
	"context"
)

import (
	"route256/libs/db/transaction"
)

type Domain interface {
	AddToCart(ctx context.Context, itemInfo *CartItem) error
	GetListItems(ctx context.Context, user int64) (*Cart, error)
	DeleteFromCart(ctx context.Context, itemInfo *CartItem) error
	Purchase(ctx context.Context, user int64) error
}

type domain struct {
	stockRepository    StockRepository
	orderRepository    OrderRepository
	transactionService transaction.Manager
	cartItemRepository CartItemRepository
	productRepository  ProductRepository
}

func New(
	stockRepository StockRepository,
	orderRepository OrderRepository,
	transactionService transaction.Manager,
	cartItemRepository CartItemRepository,
	productRepository ProductRepository,
) Domain {
	return &domain{
		stockRepository,
		orderRepository,
		transactionService,
		cartItemRepository,
		productRepository,
	}
}
