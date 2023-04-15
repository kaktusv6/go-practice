package jobs

import (
	"context"
	"time"
)

import (
	"go.uber.org/zap"
	"route256/libs/logger"
	"route256/libs/pool/batch"
	"route256/loms/internal/domain"
)

// OrdersChecker CRON для проверки заказоы
type OrdersChecker struct {
	domain domain.Domain
}

func NewOrdersChecker(
	domain domain.Domain,
) *OrdersChecker {
	return &OrdersChecker{
		domain: domain,
	}
}

func (o *OrdersChecker) Run() {
	ctx := context.Background()

	logger.Info("Get orders")

	orders, err := o.domain.GetAll(ctx)
	if err != nil {
		logger.Error("Error get orders", zap.Error(err))
		return
	}
	logger.Info("Get orders count", zap.Int("count", len(orders)))

	now := time.Now()

	tasks := make([]batch.Task[*domain.Order, *domain.Order], 0, len(orders))

	for _, order := range orders {
		tasks = append(tasks, batch.Task[*domain.Order, *domain.Order]{
			Callback: func(order *domain.Order) *domain.Order {
				logger.Info("Start checking for time payment", zap.Int64("orderId", order.ID))
				isOrderStatusValid := order.Status == domain.AwaitingPayment
				isOrderCreatedAtValid := now.Sub(order.UpdatedAt) >= (10 * time.Minute)

				if isOrderStatusValid && isOrderCreatedAtValid {
					err = o.domain.FailOrder(ctx, order)
					if err != nil {
						logger.Error("Error fail order", zap.Error(err), zap.Int64("orderId", order.ID))
					} else {
						logger.Info("Order mark as fail", zap.Int64("orderId", order.ID))
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
