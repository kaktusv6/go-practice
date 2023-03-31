package domain

import (
	"context"
	"route256/libs/db/transaction"
)

type Domain interface {
	CreateOrder(ctx context.Context, order *Order) (int64, error)
	GetListOrder(ctx context.Context, orderID int64) (*Order, error)
	OrderPayedMark(ctx context.Context, order *Order) error
	GetStocksBySKU(ctx context.Context, sku uint32) ([]*Stock, error)
	CancelOrder(ctx context.Context, order *Order) error
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
	manager                  transaction.Manager
	stockRepository          StockRepository
	orderRepository          OrderRepository
	orderItemRepository      OrderItemRepository
	orderItemStockRepository OrderItemStockRepository
}

func NewDomain(
	manager transaction.Manager,
	stockRepository StockRepository,
	orderRepository OrderRepository,
	orderItemRepository OrderItemRepository,
	orderItemStockRepository OrderItemStockRepository,
) Domain {
	return &domain{
		manager:                  manager,
		stockRepository:          stockRepository,
		orderRepository:          orderRepository,
		orderItemRepository:      orderItemRepository,
		orderItemStockRepository: orderItemStockRepository,
	}
}
