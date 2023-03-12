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
)

func (d *domain) AddToCart(ctx context.Context, itemInfo ItemInfo) error {
	result, err := d.lomsClient.Stocks(ctx, &lomsV1Clinet.StocksRequest{Sku: itemInfo.Sku})
	if err != nil {
		return errors.WithMessage(err, "get stocks from loms")
	}

	counter := int64(itemInfo.Count)
	for _, stock := range result.GetStocks() {
		counter -= int64(stock.GetCount())
		if counter <= 0 {
			return nil
		}
	}

	return ErrInsufficientStocks
}
