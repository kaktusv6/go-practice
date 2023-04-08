package tests

import (
	"context"
	"testing"
)

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"route256/libs/db"
	"route256/libs/db/transaction"
	dbMocks "route256/loms/internal/client/db/mocks"
	"route256/loms/internal/domain"
	domainMocks "route256/loms/internal/domain/mocks"
)

func TestCreateOrder(t *testing.T) {
	var (
		mc = minimock.NewController(t)

		ctx    = context.Background()
		txOpts = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
		tx     = dbMocks.NewTxMock(mc)
		ctxTx  = context.WithValue(ctx, db.TxKey, tx)

		amountItems             = 5
		amountStocksOnItem      = 2
		user                    = gofakeit.Int64()
		items                   []*domain.Item
		itemStocks              [][]*domain.Stock
		invalidaOrderItemStocks [][]*domain.Stock

		emptyOrderId        = int64(0)
		order               *domain.Order
		orderWithEmptyItems = &domain.Order{}

		saveOrderError           = errors.New("error save orderID")
		updateOrderError         = errors.New("error update orderID")
		saveOrderItemsError      = errors.New("error save orderID items")
		getStocksError           = errors.New("error get stocks")
		saveOrderItemStocksError = errors.New("error save orderID item stock")
	)

	items = make([]*domain.Item, 0, amountItems)
	itemStocks = make([][]*domain.Stock, 0, amountItems)
	invalidaOrderItemStocks = make([][]*domain.Stock, 0, amountItems)
	for i := 0; i < amountItems; i++ {
		itemCount := gofakeit.Uint16()
		items = append(items, &domain.Item{
			Sku:   gofakeit.Uint32(),
			Count: itemCount,
		})
		stocks := make([]*domain.Stock, 0, amountStocksOnItem)
		invalidStocks := make([]*domain.Stock, 0, amountStocksOnItem)
		for j := 0; j < amountStocksOnItem; j++ {
			stocks = append(stocks, &domain.Stock{
				WarehouseID: gofakeit.Int64(),
				Sku:         gofakeit.Uint32(),
				Count:       uint64(itemCount),
			})
			invalidStocks = append(invalidStocks, &domain.Stock{
				WarehouseID: gofakeit.Int64(),
				Sku:         gofakeit.Uint32(),
				Count:       uint64(itemCount / uint16(amountStocksOnItem+1)),
			})
		}
		itemStocks = append(itemStocks, stocks)
		invalidaOrderItemStocks = append(invalidaOrderItemStocks, invalidStocks)
	}

	order = &domain.Order{
		ID:    gofakeit.Int64(),
		User:  user,
		Items: items,
	}

	testCases := []struct {
		name string
		args struct {
			order *domain.Order
		}
		want                         int64
		error                        error
		dbMock                       dbMock
		stockRepositoryMock          stockRepositoryMock
		orderRepositoryMock          orderRepositoryMock
		orderItemRepositoryMock      orderItemRepositoryMock
		orderItemStockRepositoryMock orderItemStockRepositoryMock
		orderStatusNotifierMock      orderStatusNotifierMock
	}{
		{
			name: "negative case empty items",
			args: struct {
				order *domain.Order
			}{
				order: orderWithEmptyItems,
			},
			want:  emptyOrderId,
			error: domain.ErrorEmptyItems,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
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
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				return nil
			},
		},
		{
			name: "negative case error save orderID",
			args: struct {
				order *domain.Order
			}{
				order: order,
			},
			want:  emptyOrderId,
			error: saveOrderError,
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.RollbackMock.Expect(ctxTx).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.SaveMock.Expect(ctxTx, order).Return(saveOrderError)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				return nil
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				return nil
			},
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				return nil
			},
		},
		{
			name: "negative case error save orderID items",
			args: struct {
				order *domain.Order
			}{
				order: order,
			},
			want:  emptyOrderId,
			error: saveOrderItemsError,
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.RollbackMock.Expect(ctxTx).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.SaveMock.Expect(ctxTx, order).Return(nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.SaveManyMock.Expect(ctxTx, order.ID, order.Items).Return(saveOrderItemsError)
				return mock
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				return nil
			},
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				return nil
			},
		},
		{
			name: "negative case error get stocks",
			args: struct {
				order *domain.Order
			}{
				order: order,
			},
			want:  emptyOrderId,
			error: getStocksError,
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.RollbackMock.Expect(ctxTx).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySKUMock.Expect(ctxTx, items[0].Sku).Return(nil, getStocksError)
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.SaveMock.Expect(ctxTx, order).Return(nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.SaveManyMock.Expect(ctxTx, order.ID, order.Items).Return(nil)
				return mock
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				return nil
			},
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				return nil
			},
		},
		{
			name: "negative case error save orderID item stock",
			args: struct {
				order *domain.Order
			}{
				order: order,
			},
			want:  emptyOrderId,
			error: saveOrderItemStocksError,
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.RollbackMock.Expect(ctxTx).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				for index, item := range items {
					mock.GetListBySKUMock.When(ctxTx, item.Sku).Then(itemStocks[index], nil)
				}
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.SaveMock.Expect(ctxTx, order).Return(nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.SaveManyMock.Expect(ctxTx, order.ID, order.Items).Return(nil)
				return mock
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				mock := domainMocks.NewOrderItemStockRepositoryMock(mc)
				mock.SaveMock.Return(saveOrderItemStocksError)
				return mock
			},
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				return nil
			},
		},
		{
			name: "negative case error save orderID",
			args: struct {
				order *domain.Order
			}{
				order: order,
			},
			want:  emptyOrderId,
			error: updateOrderError,
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.RollbackMock.Expect(ctxTx).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				for index, item := range items {
					mock.GetListBySKUMock.When(ctxTx, item.Sku).Then(invalidaOrderItemStocks[index], nil)
				}
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.SaveMock.Expect(ctxTx, order).Return(nil)
				mock.UpdateMock.Expect(ctxTx, order).Return(updateOrderError)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.SaveManyMock.Expect(ctxTx, order.ID, order.Items).Return(nil)
				return mock
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				mock := domainMocks.NewOrderItemStockRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				return nil
			},
		},
		{
			name: "positive case orderID create with status failed",
			args: struct {
				order *domain.Order
			}{
				order: order,
			},
			want:  order.ID,
			error: nil,
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctxTx).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				for index, item := range items {
					mock.GetListBySKUMock.When(ctxTx, item.Sku).Then(invalidaOrderItemStocks[index], nil)
				}
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.SaveMock.Expect(ctxTx, order).Return(nil)
				mock.UpdateMock.Expect(ctxTx, order).Return(nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.SaveManyMock.Expect(ctxTx, order.ID, order.Items).Return(nil)
				return mock
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				mock := domainMocks.NewOrderItemStockRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				mock := domainMocks.NewOrderStatusNotifierMock(mc)
				mock.NotifyMock.Return(nil)
				return mock
			},
		},
		{
			name: "positive case",
			args: struct {
				order *domain.Order
			}{
				order: order,
			},
			want:  order.ID,
			error: nil,
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctxTx).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				for index, item := range items {
					mock.GetListBySKUMock.When(ctxTx, item.Sku).Then(itemStocks[index], nil)
				}
				return mock
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.SaveMock.Expect(ctxTx, order).Return(nil)
				mock.UpdateMock.Expect(ctxTx, order).Return(nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.SaveManyMock.Expect(ctxTx, order.ID, order.Items).Return(nil)
				return mock
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				mock := domainMocks.NewOrderItemStockRepositoryMock(mc)
				mock.SaveMock.Return(nil)
				return mock
			},
			orderStatusNotifierMock: func(mc *minimock.Controller) domain.OrderStatusNotifier {
				mock := domainMocks.NewOrderStatusNotifierMock(mc)
				mock.NotifyMock.Return(nil)
				return mock
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
				testCase.orderStatusNotifierMock(mc),
			)

			res, err := d.CreateOrder(ctx, testCase.args.order)
			require.Equal(t, testCase.want, res)
			if testCase.error == nil {
				require.Equal(t, testCase.error, err)
			} else {
				require.ErrorContains(t, err, testCase.error.Error())
			}
		})
	}
}
