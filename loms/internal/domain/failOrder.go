package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
)

var (
	ErrorOrderIsNew = errors.New("order is new")
)

func (d *domain) FailOrder(ctx context.Context, order *Order) error {
	if order.Status == New {
		return ErrorOrderIsNew
	}

	if order.Status == Cancelled {
		return ErrorOrderIsCanceled
	}

	if order.Status == Failed {
		return ErrorOrderIsFailed
	}

	err := d.transactionManager.RunRepeatableReade(ctx, func(ctxTx context.Context) error {
		isRevertStocks := order.Status == Payed

		order.Status = Failed

		err := d.orderRepository.Update(ctxTx, order)
		if err != nil {
			return err
		}

		if isRevertStocks {
			err = d.revertOrderItemsToWarehouses(ctxTx, order)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
