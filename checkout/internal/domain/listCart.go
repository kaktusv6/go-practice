package domain

import (
	"context"
)

import (
	productServiceV1Clinet "route256/checkout/pkg/product_service_v1"
)

func (d *domain) GetListItems(ctx context.Context, user int64) (*Cart, error) {
	userCartItems, err := d.cartItemRepository.GetUserCartItems(ctx, user)
	if err != nil {
		return nil, err
	}

	cart := &Cart{}
	for index, cartItem := range userCartItems {
		product, err := d.productServiceClient.GetProduct(ctx, &productServiceV1Clinet.GetProductRequest{
			Token: d.productServiceToken,
			Sku:   cartItem.Sku,
		})
		if err != nil {
			return cart, err
		}
		userCartItems[index].Product = &ProductInfo{
			Name:  product.GetName(),
			Price: product.GetPrice(),
		}
	}

	cart.Items = userCartItems
	cart.calculateTotalPrice()

	return cart, nil
}
