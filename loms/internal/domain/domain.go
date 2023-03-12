package domain

import "context"

type Domain interface {
	CreateOrder(ctx context.Context, user int64, items []Item) (int64, error)
	GetListOrder(ctx context.Context, orderID int64) (Order, error)
	OrderPayedMark(ctx context.Context, orderID int64) error
	GetStocksBySKU(ctx context.Context, sku uint32) ([]Stock, error)
	CancelOrder(ctx context.Context, orderID int64) error
}

const (
	New             = "new"
	AwaitingPayment = "awaiting payment"
	Failed          = "failed"
	Payed           = "payed"
	Cancelled       = "cancelled"
)

type domain struct {
}

func NewDomain() Domain {
	return &domain{}
}
