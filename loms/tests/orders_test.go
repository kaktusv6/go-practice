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

func TestOrders(t *testing.T) {
	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		amountOrders = 5

		orderId = gofakeit.Int64()
		orders  []*domain.Order

		getOrderError = errors.New("error get orders")
	)

	orders = make([]*domain.Order, 0, amountOrders)
	for i := 0; i < amountOrders; i++ {
		orders = append(orders, &domain.Order{
			ID: orderId,
			Status: gofakeit.RandomString([]string{
				domain.New,
				domain.AwaitingPayment,
				domain.Failed,
				domain.Payed,
				domain.Cancelled,
			}),
			User:      gofakeit.Int64(),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		})
	}

	testCases := []struct {
		name                         string
		want                         []*domain.Order
		error                        error
		dbMock                       dbMock
		stockRepositoryMock          stockRepositoryMock
		orderRepositoryMock          orderRepositoryMock
		orderItemRepositoryMock      orderItemRepositoryMock
		orderItemStockRepositoryMock orderItemStockRepositoryMock
	}{
		{
			name:  "positive case",
			want:  orders,
			error: nil,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.GetAllMock.Expect(ctx).Return(orders, nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				return nil
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				return nil
			},
		},
		{
			name:  "negative case error get orders",
			want:  nil,
			error: getOrderError,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.GetAllMock.Expect(ctx).Return(nil, getOrderError)
				return mock
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
			)

			res, err := d.GetAll(ctx)
			require.Equal(t, testCase.want, res)
			require.Equal(t, testCase.error, err)
		})
	}
}
