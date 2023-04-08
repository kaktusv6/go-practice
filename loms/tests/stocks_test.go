package tests

import (
	"context"
	"testing"
)

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"route256/libs/db"
	"route256/libs/db/transaction"
	"route256/loms/internal/domain"
	domainMocks "route256/loms/internal/domain/mocks"
)

func TestStocks(t *testing.T) {
	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		amountStocks = 5

		sku    = gofakeit.Uint32()
		stocks []*domain.Stock

		getStocksError = errors.New("error get orderID items by orderID id")
	)

	stocks = make([]*domain.Stock, 0, amountStocks)
	for i := 0; i < amountStocks; i++ {
		stocks = append(stocks, &domain.Stock{
			WarehouseID: gofakeit.Int64(),
			Sku:         gofakeit.Uint32(),
			Count:       gofakeit.Uint64(),
		})
	}

	testCases := []struct {
		name                         string
		want                         []*domain.Stock
		error                        error
		dbMock                       dbMock
		stockRepositoryMock          stockRepositoryMock
		orderRepositoryMock          orderRepositoryMock
		orderItemRepositoryMock      orderItemRepositoryMock
		orderItemStockRepositoryMock orderItemStockRepositoryMock
	}{
		{
			name:  "positive case",
			want:  stocks,
			error: nil,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySKUMock.Expect(ctx, sku).Return(stocks, nil)
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				return nil
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				return nil
			},
		},
		{
			name:  "negative case error get orderID",
			want:  nil,
			error: getStocksError,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySKUMock.Expect(ctx, sku).Return(nil, getStocksError)
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				return nil
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				return nil
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			d := domain.NewDomain(
				transaction.NewTransactionManager(testCase.dbMock(mc)),
				testCase.stockRepositoryMock(mc),
				testCase.orderRepositoryMock(mc),
				testCase.orderItemRepositoryMock(mc),
				testCase.orderItemStockRepositoryMock(mc),
				nil,
			)

			res, err := d.GetStocksBySKU(ctx, sku)
			require.Equal(t, testCase.want, res)
			require.Equal(t, testCase.error, err)
		})
	}
}
