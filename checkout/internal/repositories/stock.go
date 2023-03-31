package repositories

import (
	"context"
)

import (
	"route256/checkout/internal/domain"
	lomsV1Clinet "route256/loms/pkg/loms_v1"
)

type StockRepository struct {
	lomsClient lomsV1Clinet.LomsV1Client
}

func NewStockRepository(lomsClient lomsV1Clinet.LomsV1Client) domain.StockRepository {
	return &StockRepository{
		lomsClient: lomsClient,
	}
}

func (s *StockRepository) GetListBySku(ctx context.Context, sku uint32) ([]*domain.Stock, error) {
	stockResponse, err := s.lomsClient.Stocks(ctx, &lomsV1Clinet.StocksRequest{Sku: sku})
	if err != nil {
		return nil, err
	}

	stocks := make([]*domain.Stock, 0, len(stockResponse.GetStocks()))
	for _, stock := range stockResponse.GetStocks() {
		stocks = append(stocks, &domain.Stock{
			WarehouseID: stock.GetWarehouseId(),
			Count:       stock.GetCount(),
		})
	}

	return stocks, nil
}
