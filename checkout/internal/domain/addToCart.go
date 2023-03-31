package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrProductNotFound    = errors.New("product not found by sku")
)

func (d *domain) AddToCart(ctx context.Context, cartItem *CartItem) error {
	existCartItem, err := d.cartItemRepository.GetOne(ctx, cartItem.User, cartItem.Sku)
	if err != nil {
		return err
	}

	stocks, err := d.stockRepository.GetListBySku(ctx, cartItem.Sku)
	if err != nil {
		return err
	}

	count := int64(cartItem.Count)
	for _, stock := range stocks {
		count -= int64(stock.Count)
	}

	if count > 0 {
		return ErrInsufficientStocks
	}

	isCreateCartItem := existCartItem == nil

	if isCreateCartItem {
		productInfo, err := d.productRepository.GetProductBySku(ctx, cartItem.Sku)
		if err != nil {
			return err
		}
		if productInfo == nil {
			return ErrProductNotFound
		}
	} else {
		existCartItem.Count = existCartItem.Count + cartItem.Count
	}

	err = d.transactionService.RepeatableRead(ctx, func(ctxTx context.Context) error {
		if isCreateCartItem {
			err = d.cartItemRepository.Create(ctxTx, cartItem)
		} else {
			err = d.cartItemRepository.Update(ctxTx, existCartItem)
		}
		return err
	})

	return err
}
