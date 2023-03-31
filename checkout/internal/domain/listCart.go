package domain

import (
	"context"
)

func (d *domain) GetListItems(ctx context.Context, user int64) (*Cart, error) {
	userCartItems, err := d.cartItemRepository.GetUserCartItems(ctx, user)
	if err != nil {
		return nil, err
	}

	cart := &Cart{}

	skus := make([]uint32, 0, len(userCartItems))
	for _, userCartItem := range userCartItems {
		skus = append(skus, userCartItem.Sku)
	}

	productInfoList, err := d.productRepository.GetListBySkus(ctx, skus)
	if err != nil {
		return nil, err
	}
	for index, productInfo := range productInfoList {
		userCartItems[index].Product = productInfo
	}

	cart.Items = userCartItems
	cart.CalculateTotalPrice()

	return cart, nil
}
