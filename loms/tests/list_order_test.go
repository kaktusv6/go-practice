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

func TestListOrder(t *testing.T) {
	var (
		mc  = minimock.NewController(t)
		ctx = context.Background()

		amountOrderItems = 5

		orderId    = gofakeit.Int64()
		order      *domain.Order
		emptyOrder *domain.Order
		orderItems []*domain.Item

		getOrderError      = errors.New("error get orderID by id")
		getOrderItemsError = errors.New("error get orderID items by orderID id")
	)

	orderItems = make([]*domain.Item, 0, amountOrderItems)
	for i := 0; i < amountOrderItems; i++ {
		orderItems = append(orderItems, &domain.Item{
			Sku:   gofakeit.Uint32(),
			Count: gofakeit.Uint16(),
		})
	}

	emptyOrder = &domain.Order{}
	order = &domain.Order{
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
	}

	testCases := []struct {
		name                         string
		want                         *domain.Order
		error                        error
		dbMock                       dbMock
		stockRepositoryMock          stockRepositoryMock
		orderRepositoryMock          orderRepositoryMock
		orderItemRepositoryMock      orderItemRepositoryMock
		orderItemStockRepositoryMock orderItemStockRepositoryMock
	}{
		{
			name:  "positive case",
			want:  order,
			error: nil,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.GetByIdMock.Expect(ctx, orderId).Return(order, nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.GetByOrderIdMock.Expect(ctx, orderId).Return(orderItems, nil)
				return mock
			},
			orderItemStockRepositoryMock: func(mc *minimock.Controller) domain.OrderItemStockRepository {
				return nil
			},
		},
		{
			name:  "negative case error get orderID",
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
				mock.GetByIdMock.Expect(ctx, orderId).Return(nil, getOrderError)
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
			name:  "negative case no orderID by id",
			want:  nil,
			error: domain.ErrorOrderNotFound,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.GetByIdMock.Expect(ctx, orderId).Return(emptyOrder, nil)
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
			name:  "negative case error get orderID items",
			want:  nil,
			error: getOrderItemsError,
			dbMock: func(mc *minimock.Controller) db.DB {
				return nil
			},
			stockRepositoryMock: func(mc *minimock.Controller) domain.StockRepository {
				return nil
			},
			orderRepositoryMock: func(mc *minimock.Controller) domain.OrderRepository {
				mock := domainMocks.NewOrderRepositoryMock(mc)
				mock.GetByIdMock.Expect(ctx, orderId).Return(order, nil)
				return mock
			},
			orderItemRepositoryMock: func(mc *minimock.Controller) domain.OrderItemRepository {
				mock := domainMocks.NewOrderItemRepositoryMock(mc)
				mock.GetByOrderIdMock.Expect(ctx, orderId).Return(nil, getOrderItemsError)
				return mock
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

			res, err := d.GetListOrder(ctx, orderId)
			require.Equal(t, testCase.want, res)
			require.Equal(t, testCase.error, err)
		})
	}
}
