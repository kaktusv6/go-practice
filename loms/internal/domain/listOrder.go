package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
)

var (
	ErrorOrderNotFound = errors.New("order not found")
)

func (d *domain) GetListOrder(ctx context.Context, orderID int64) (*Order, error) {
	order, err := d.orderRepository.GetById(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.ID == 0 {
		return nil, ErrorOrderNotFound
	}

	order.Items, err = d.orderItemRepository.GetByOrderId(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
