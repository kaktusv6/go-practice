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

func TestDeleteFromCart(t *testing.T) {
	var (
		mc    = minimock.NewController(t)
		ctx   = context.Background()
		tx    = dbMocks.NewTxMock(t)
		ctxTx = context.WithValue(ctx, db.TxKey, tx)

		cartItem                        *domain.CartItem
		existCartItem                   *domain.CartItem
		existCartItemWithIncrimentCount *domain.CartItem

		txOpts = pgx.TxOptions{IsoLevel: pgx.RepeatableRead}

		getExistCartItemError = errors.New("error get exist cart item")
		deleteCartItemError   = errors.New("error delete item error")
		updateCartItemError   = errors.New("error update item error")
	)

	fakeCartItem := domain.CartItem{
		User:  gofakeit.Int64(),
		Sku:   gofakeit.Uint32(),
		Count: gofakeit.Uint16(),
	}
	cartItem = &fakeCartItem
	existCartItem = &fakeCartItem
	existCartItemWithIncrimentCount = &domain.CartItem{
		User:  gofakeit.Int64(),
		Sku:   gofakeit.Uint32(),
		Count: cartItem.Count + 1,
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
			name: "positive case with delete cart item",
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
				mock.DeleteMock.Expect(ctxTx, existCartItem).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
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
			name: "positive case with update cart item",
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
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(existCartItemWithIncrimentCount, nil)
				mock.UpdateMock.Expect(ctxTx, existCartItemWithIncrimentCount).Return(nil)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
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
			name: "negative case error delete cart item",
			args: struct {
				cartItem *domain.CartItem
			}{
				cartItem: cartItem,
			},
			want: deleteCartItemError,
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := domainMocks.NewCartItemRepositoryMock(mc)
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(existCartItem, nil)
				mock.DeleteMock.Expect(ctxTx, cartItem).Return(deleteCartItemError)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
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
				mock.GetOneMock.Expect(ctx, cartItem.User, cartItem.Sku).Return(existCartItemWithIncrimentCount, nil)
				mock.UpdateMock.Expect(ctxTx, existCartItemWithIncrimentCount).Return(updateCartItemError)
				return mock
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
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
			res := d.DeleteFromCart(ctx, testCase.args.cartItem)
			if testCase.want == nil {
				require.Equal(t, testCase.want, res)
			} else {
				require.ErrorContains(t, res, testCase.want.Error())
			}
		})
	}
}
