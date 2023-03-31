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
	dbMocks "route256/checkout/internal/client/db/mocks"
	"route256/checkout/internal/domain"
	domainMocks "route256/checkout/internal/domain/mocks"
	"route256/libs/db"
	"route256/libs/db/transaction"
)

func TestAddToCart(t *testing.T) {
	var (
		mc    = minimock.NewController(t)
		ctx   = context.Background()
		tx    = dbMocks.NewTxMock(t)
		ctxTx = context.WithValue(ctx, db.TxKey, tx)

		amountStocks = uint16(2)

		cartItem      *domain.CartItem
		existCartItem *domain.CartItem
		stocks        []*domain.Stock
		invalidStocks []*domain.Stock
		txOpts        = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
		productInfo   *domain.ProductInfo

		getExistCartItemError = errors.New("error get exist cart item")
		getStocksError        = errors.New("error get stocks")
		getProductInfoError   = errors.New("error get product info")
		createCartItemError   = errors.New("error create item error")
		updateCartItemError   = errors.New("error update item error")
	)

	fakeCartItem := domain.CartItem{
		User:  gofakeit.Int64(),
		Sku:   gofakeit.Uint32(),
		Count: gofakeit.Uint16() * amountStocks,
	}
	cartItem = &fakeCartItem
	existCartItem = &fakeCartItem

	stocks = make([]*domain.Stock, 0, amountStocks)
	for i := 0; i < int(amountStocks); i++ {
		stocks = append(stocks, &domain.Stock{
			WarehouseID: gofakeit.Int64(),
			Count:       uint64(cartItem.Count),
		})
	}
	invalidStocks = make([]*domain.Stock, 0, amountStocks)
	invalidStocks = append(invalidStocks, &domain.Stock{
		WarehouseID: gofakeit.Int64(),
		Count:       uint64(cartItem.Count / amountStocks),
	})
	productInfo = &domain.ProductInfo{
		Name:  gofakeit.BeerName(),
		Price: uint32(gofakeit.Price(100, 5000)),
	}

	testCases := []struct {
		name string
		args struct {
			cartItem *domain.CartItem
		}
		want                   error
		stockRepositoryMock    stockRepositoryMock
		orderRepositoryMock    orderRepositoryMock
		cartItemRepositoryMock cartItemRepositoryMock
		productRepositoryMock  productRepositoryMock
		dbMock                 dbMock
	}{
		{
			name: "positive case with exist item to cart",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: nil,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(existCartItem, nil)
				mock.UpdateMock.Expect(ctxTx, existCartItem).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySkuMock.Expect(ctx, cartItem.Sku).Return(stocks, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				return nil
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctxTx).Return(nil)
				return mock
			},
		},
		{
			name: "positive case with new item to cart",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: nil,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(nil, nil)
				mock.CreateMock.Expect(ctxTx, cartItem).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySkuMock.Expect(ctx, cartItem.Sku).Return(stocks, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				mock := domainMocks.NewProductRepositoryMock(mc)
				mock.GetProductBySkuMock.Expect(ctx, cartItem.Sku).Return(productInfo, nil)
				return mock
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.CommitMock.Expect(ctxTx).Return(nil)
				return mock
			},
		},
		{
			name: "negative case error get exist cart item",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: getExistCartItemError,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(nil, getExistCartItemError)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				return nil
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
		},
		{
			name: "negative case error get stocks",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: getStocksError,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(nil, nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySkuMock.Expect(ctx, cartItem.Sku).Return(nil, getStocksError)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				return nil
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
		},
		{
			name: "negative case error insufficient stocks",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: domain.ErrInsufficientStocks,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(nil, nil)
				mock.CreateMock.Expect(ctxTx, cartItem).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySkuMock.Expect(ctx, cartItem.Sku).Return(invalidStocks, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				return nil
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
		},
		{
			name: "negative case error get product info",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: getProductInfoError,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(nil, nil)
				mock.CreateMock.Expect(ctxTx, cartItem).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySkuMock.Expect(ctx, cartItem.Sku).Return(stocks, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				mock := domainMocks.NewProductRepositoryMock(mc)
				mock.GetProductBySkuMock.Expect(ctx, cartItem.Sku).Return(nil, getProductInfoError)
				return mock
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
		},
		{
			name: "negative case error create cart item",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: createCartItemError,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(nil, nil)
				mock.CreateMock.Expect(ctxTx, cartItem).Return(createCartItemError)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySkuMock.Expect(ctx, cartItem.Sku).Return(stocks, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				mock := domainMocks.NewProductRepositoryMock(mc)
				mock.GetProductBySkuMock.Expect(ctx, cartItem.Sku).Return(productInfo, nil)
				return mock
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.RollbackMock.Expect(ctxTx).Return(nil)
				return mock
			},
		},
		{
			name: "negative case error update cart item",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: updateCartItemError,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(existCartItem, nil)
				mock.UpdateMock.Expect(ctxTx, existCartItem).Return(updateCartItemError)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				mock := domainMocks.NewStockRepositoryMock(mc)
				mock.GetListBySkuMock.Expect(ctx, cartItem.Sku).Return(stocks, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				return nil
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				mock := dbMocks.NewDBMock(mc)
				mock.BeginTxMock.Expect(ctx, txOpts).Return(tx, nil)
				tx.RollbackMock.Expect(ctxTx).Return(nil)
				return mock
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			d := domain.New(
				testCase.stockRepositoryMock(mc),
				testCase.orderRepositoryMock(mc),
				transaction.NewTransactionManager(testCase.dbMock(mc)),
				testCase.cartItemRepositoryMock(mc),
				testCase.productRepositoryMock(mc),
			)
			res := d.AddToCart(ctx, testCase.args.cartItem)
			if testCase.want == nil {
				require.Equal(t, testCase.want, res)
			} else {
				require.ErrorContains(t, res, testCase.want.Error())
			}
		})
	}
}
