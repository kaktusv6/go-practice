package lomsV1

import (
	"context"
)

import (
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) Stocks(ctx context.Context, req *desc.StocksRequest) (*desc.StocksResponse, error) {
	stocks, err := i.domain.GetStocksBySKU(ctx, req.GetSku())
	if err != nil {
		return nil, err
	}

	resStocks := make([]*desc.Stock, 0, len(stocks))
	for _, stock := range stocks {
		resStocks = append(resStocks, &desc.Stock{
			Count:       stock.Count,
			WarehouseId: stock.WarehouseID,
		})
	}

	return &desc.StocksResponse{
		Stocks: resStocks,
	}, nil
}
