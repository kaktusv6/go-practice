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
	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/mocks"
	"route256/libs/db"
	"route256/libs/db/transaction"
)

func TestPurchase(t *testing.T) {
	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		user            = gofakeit.Int64()
		amountCartItems = 5

		cartItems    []*domain.CartItem
		cartItemSkus []uint32
		productItems []*domain.ProductInfo

		cartItemRepoError = errors.New("Error get user cart items")
		orderRepoError    = errors.New("Error create order")
	)

	cartItems = make([]*domain.CartItem, 0, amountCartItems)
	for i := 0; i < amountCartItems; i++ {
		cartItems = append(cartItems, &domain.CartItem{
			User:  gofakeit.Int64(),
			Sku:   gofakeit.Uint32(),
			Count: gofakeit.Uint16(),
		})
	}

	testCases := []struct {
		name string
		args struct {
			ctx  context.Context
			user int64
		}
		want                   error
		stockRepositoryMock    stockRepositoryMock
		orderRepositoryMock    orderRepositoryMock
		cartItemRepositoryMock cartItemRepositoryMock
		productRepositoryMock  productRepositoryMock
		dbMock                 dbMock
	}{
		{
			name: "positive case",
			args: struct {
				ctx  context.Context
				user int64
			}{
				ctx:  ctx,
				user: user,
			},
			want: nil,
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := mocks.NewOrderRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, user, cartItems).Return(nil)
				return mock
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := mocks.NewCartItemRepositoryMock(mc)
				mock.GetUserCartItemsMock.Expect(ctx, user).Return(cartItems, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				mock := mocks.NewProductRepositoryMock(mc)
				mock.GetListBySkusMock.Expect(ctx, cartItemSkus).Return(productItems, nil)
				return mock
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
		},
		{
			name: "negative case error get cart items",
			args: struct {
				ctx  context.Context
				user int64
			}{
				ctx:  ctx,
				user: user,
			},
			want: cartItemRepoError,
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := mocks.NewCartItemRepositoryMock(mc)
				mock.GetUserCartItemsMock.Expect(ctx, user).Return(nil, cartItemRepoError)
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
			name: "negative case empty cart items",
			args: struct {
				ctx  context.Context
				user int64
			}{
				ctx:  ctx,
				user: user,
			},
			want: domain.CartItemsEmptyError,
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := mocks.NewCartItemRepositoryMock(mc)
				mock.GetUserCartItemsMock.Expect(ctx, user).Return(make([]*domain.CartItem, 0), nil)
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
			name: "negative case by order create",
			args: struct {
				ctx  context.Context
				user int64
			}{
				ctx:  ctx,
				user: user,
			},
			want: orderRepoError,
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := mocks.NewOrderRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, user, cartItems).Return(orderRepoError)
				return mock
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := mocks.NewCartItemRepositoryMock(mc)
				mock.GetUserCartItemsMock.Expect(ctx, user).Return(cartItems, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				return nil
			},
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
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
			res := d.Purchase(testCase.args.ctx, testCase.args.user)
			require.Equal(t, testCase.want, res)
		})
	}
}
