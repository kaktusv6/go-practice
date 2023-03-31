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

func TestGetListItems(t *testing.T) {

	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		user            = gofakeit.Int64()
		amountCartItems = 5

		cartResult        *domain.Cart
		cartItems         []*domain.CartItem
		cartItemSkus      []uint32
		productItems      []*domain.ProductInfo
		cartItemRepoError = errors.New("Error get user cart items")
		productRepoError  = errors.New("Error get product info")
	)

	cartItems = make([]*domain.CartItem, 0, amountCartItems)
	productItems = make([]*domain.ProductInfo, 0, amountCartItems)
	for i := 0; i < amountCartItems; i++ {
		product := &domain.ProductInfo{
			Name:  gofakeit.BeerName(),
			Price: uint32(gofakeit.Price(10, 5000)),
		}
		item := &domain.CartItem{
			User:    gofakeit.Int64(),
			Sku:     gofakeit.Uint32(),
			Count:   gofakeit.Uint16(),
			Product: product,
		}
		cartItems = append(cartItems, item)
		cartItemSkus = append(cartItemSkus, item.Sku)
		productItems = append(productItems, product)
	}

	cartResult = &domain.Cart{
		Items: cartItems,
	}
	cartResult.CalculateTotalPrice()

	testCases := []struct {
		name string
		args struct {
			user int64
		}
		want                   *domain.Cart
		error                  error
		stockRepositoryMock    stockRepositoryMock
		orderRepositoryMock    orderRepositoryMock
		cartItemRepositoryMock cartItemRepositoryMock
		productRepositoryMock  productRepositoryMock
		dbMock                 dbMock
	}{
		{
			name: "positive case",
			args: struct {
				user int64
			}{
				user: user,
			},
			want:  cartResult,
			error: nil,
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
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
			name: "negative case by cart items",
			args: struct {
				user int64
			}{
				user: user,
			},
			want:  nil,
			error: cartItemRepoError,
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
			name: "negative case by products",
			args: struct {
				user int64
			}{
				user: user,
			},
			want:  nil,
			error: productRepoError,
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				return nil
			},
			cartItemRepositoryMock: func(mc *minimock.Controller) domain.CartItemRepository {
				mock := mocks.NewCartItemRepositoryMock(mc)
				mock.GetUserCartItemsMock.Expect(ctx, user).Return(cartItems, nil)
				return mock
			},
			productRepositoryMock: func(mc *minimock.Controller) domain.ProductRepository {
				mock := mocks.NewProductRepositoryMock(mc)
				mock.GetListBySkusMock.Expect(ctx, cartItemSkus).Return(nil, productRepoError)
				return mock
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
			res, err := d.GetListItems(ctx, testCase.args.user)
			require.Equal(t, testCase.want, res)
			require.Equal(t, testCase.error, err)
		})
	}
}
