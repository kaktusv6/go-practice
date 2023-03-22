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

	productInfoList, _ := d.productRepository.GetListBySkus(ctx, skus)
	for index, productInfo := range productInfoList {
		userCartItems[index].Product = productInfo
	}

	cart.Items = userCartItems
	cart.calculateTotalPrice()

	return cart, nil
}
