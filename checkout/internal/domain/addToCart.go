package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
	lomsV1Clinet "route256/loms/pkg/loms_v1"
)

var (
	ErrInsufficientStocks = errors.New("insufficient stocks")
	ErrProductNotFound    = errors.New("product not found by sku")
)

func (d *domain) AddToCart(ctx context.Context, cartItem CartItem) error {
	existCartItem, err := d.cartItemRepository.GetOne(ctx, cartItem.User, cartItem.Sku)
	if err != nil {
		return err
	}

	result, err := d.lomsClient.Stocks(ctx, &lomsV1Clinet.StocksRequest{Sku: cartItem.Sku})
	if err != nil {
		return errors.WithMessage(err, "get stocks from loms")
	}

	count := int64(cartItem.Count)
	for _, stock := range result.GetStocks() {
		count -= int64(stock.GetCount())
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

	err = d.transactionManager.RunRepeatableReade(ctx, func(ctxTx context.Context) error {
		if isCreateCartItem {
			err = d.cartItemRepository.Create(ctxTx, &cartItem)
		} else {
			err = d.cartItemRepository.Update(ctxTx, existCartItem)
		}
		return err
	})

	return err
}
