package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
)

var (
	ErrorOrderAlreadyCanceled = errors.New("order already mark as cancelled")
)

func (d *domain) CancelOrder(ctx context.Context, orderID int64) error {
	order, err := d.GetListOrder(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status == Cancelled {
		return ErrorOrderAlreadyCanceled
	}

	if order.Status == Failed {
		return ErrorOrderIsFailed
	}

	err = d.transactionManager.RunRepeatableReade(ctx, func(ctxTx context.Context) error {
		isRevertStocks := order.Status == Payed || order.Status == AwaitingPayment
		order.Status = Cancelled

		err = d.orderRepository.Update(ctxTx, order)
		if err != nil {
			return err
		}

		if isRevertStocks {
			orderItemStocks, err := d.orderItemStockRepository.GetListByOrderID(ctxTx, order.ID)
			if err != nil {
				return err
			}

			for _, orderItemStock := range orderItemStocks {
				err = d.revertOrderItem(ctxTx, orderItemStock)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (d *domain) revertOrderItem(ctx context.Context, orderItemStock OrderItemStock) error {
	stock, err := d.stockRepository.GetByWarehouseIDAndSku(ctx, orderItemStock.WarehouseID, orderItemStock.Sku)
	if err != nil {
		return err
	}

	stock.Count = stock.Count + orderItemStock.Count
	err = d.stockRepository.UpdateCount(ctx, stock)
	if err != nil {
		return err
	}

	err = d.orderItemStockRepository.Delete(ctx, &orderItemStock)
	if err != nil {
		return err
	}

	return nil
}
