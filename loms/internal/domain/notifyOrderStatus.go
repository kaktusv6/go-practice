package domain

import "context"

func (d *domain) NotifyOrderStatus(ctx context.Context, order *Order) error {
	return d.orderStatusNotifier.Notify(order)
}
