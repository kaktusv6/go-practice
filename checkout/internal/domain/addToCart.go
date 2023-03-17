package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
	productServiceV1Clinet "route256/checkout/pkg/product_service_v1"
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
		err = d.checkProductBySku(ctx, cartItem.Sku)
		if err != nil {
			return err
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

func (d *domain) checkProductBySku(ctx context.Context, sku uint32) error {
	productResponse, err := d.productServiceClient.GetProduct(ctx, &productServiceV1Clinet.GetProductRequest{
		Token: d.productServiceToken,
		Sku:   sku,
	})
	if err != nil {
		return err
	}

	if productResponse == nil {
		return ErrProductNotFound
	}

	return nil
}
