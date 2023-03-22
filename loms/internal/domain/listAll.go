package domain

import "context"

func (d *domain) GetAll(ctx context.Context) ([]*Order, error) {
	return d.orderRepository.GetAll(ctx)
}
