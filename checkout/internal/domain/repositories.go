package domain

//go:generate bash -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i CartItemRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i ProductRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i StockRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i OrderRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
)

type CartItemRepository interface {
	GetUserCartItems(ctx context.Context, user int64) ([]*CartItem, error)
	GetOne(ctx context.Context, user int64, sku uint32) (*CartItem, error)
	Create(ctx context.Context, cartItem *CartItem) error
	Update(ctx context.Context, cartItem *CartItem) error
	Delete(ctx context.Context, cartItem *CartItem) error
}

type ProductRepository interface {
	GetListBySkus(ctx context.Context, skus []uint32) ([]*ProductInfo, error)
	GetProductBySku(ctx context.Context, sku uint32) (*ProductInfo, error)
}

type StockRepository interface {
	GetListBySku(ctx context.Context, sku uint32) ([]*Stock, error)
}

type OrderRepository interface {
	Create(ctx context.Context, user int64, items []*CartItem) error
}
