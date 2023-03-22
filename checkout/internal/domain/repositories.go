package domain

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
