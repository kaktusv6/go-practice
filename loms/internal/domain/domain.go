package domain

import (
	"context"
)

import (
	"route256/libs/transactor"
)

type Domain interface {
	CreateOrder(ctx context.Context, user int64, items []Item) (int64, error)
	GetListOrder(ctx context.Context, orderID int64) (*Order, error)
	OrderPayedMark(ctx context.Context, orderID int64) error
	GetStocksBySKU(ctx context.Context, sku uint32) ([]Stock, error)
	CancelOrder(ctx context.Context, orderID int64) error
	GetAll(ctx context.Context) ([]*Order, error)
	FailOrder(ctx context.Context, order *Order) error
}

const (
	New             = "new"
	AwaitingPayment = "awaiting_payment"
	Failed          = "failed"
	Payed           = "payed"
	Cancelled       = "cancelled"
)

type domain struct {
	transactionManager       *transactor.TransactionManager
	stockRepository          StockRepository
	orderRepository          OrderRepository
	orderItemRepository      OrderItemRepository
	orderItemStockRepository OrderItemStockRepository
}

func NewDomain(
	transactionManager *transactor.TransactionManager,
	stockRepository StockRepository,
	orderRepository OrderRepository,
	orderItemRepository OrderItemRepository,
	orderItemStockRepository OrderItemStockRepository,
) Domain {
	return &domain{
		transactionManager:       transactionManager,
		stockRepository:          stockRepository,
		orderRepository:          orderRepository,
		orderItemRepository:      orderItemRepository,
		orderItemStockRepository: orderItemStockRepository,
	}
}
