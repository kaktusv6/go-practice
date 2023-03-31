package jobs

import (
	"context"
	"time"
)

import (
	"github.com/robfig/cron/v3"
	"route256/libs/pool/batch"
	"route256/loms/internal/domain"
)

// OrdersChecker CRON для проверки заказоы
type OrdersChecker struct {
	domain domain.Domain
	logger cron.Logger
}

func NewOrdersChecker(
	domain domain.Domain,
	logger cron.Logger,
) *OrdersChecker {
	return &OrdersChecker{
		domain: domain,
		logger: logger,
	}
}

func (o *OrdersChecker) Run() {
	ctx := context.Background()

	o.logger.Info("Get orders")

	orders, err := o.domain.GetAll(ctx)
	if err != nil {
		o.logger.Error(err, "Error get orders")
		return
	}
	o.logger.Info("Get orders count", "count", len(orders))

	now := time.Now()

	tasks := make([]batch.Task[*domain.Order, *domain.Order], 0, len(orders))

	for _, order := range orders {
		tasks = append(tasks, batch.Task[*domain.Order, *domain.Order]{
			Callback: func(order *domain.Order) *domain.Order {
				o.logger.Info("Start checking for time payment", "orderId", order.ID)
				isOrderStatusValid := order.Status == domain.AwaitingPayment
				isOrderCreatedAtValid := now.Sub(order.UpdatedAt) >= (10 * time.Minute)

				if isOrderStatusValid && isOrderCreatedAtValid {
					err = o.domain.FailOrder(ctx, order)
					if err != nil {
						o.logger.Error(err, "Error fail order", "orderId", order.ID)
					} else {
						o.logger.Info("Order mark as fail", "orderId", order.ID)
					}
				}

				return order
			},
			InputArgs: order,
		})
	}

	pool := batch.NewPool[*domain.Order, *domain.Order](ctx, 5)

	pool.Submit(ctx, tasks)
}
