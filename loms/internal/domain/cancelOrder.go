package domain

import (
	"context"
)

func (d *domain) CancelOrder(ctx context.Context, order *Order) error {
	if order.Status == Cancelled {
		return ErrorOrderIsCanceled
	}

	if order.Status == Failed {
		return ErrorOrderIsFailed
	}

	err := d.manager.RepeatableRead(ctx, func(ctxTx context.Context) error {
		isRevertStocks := order.Status == Payed
		order.Status = Cancelled

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

func (d *domain) revertOrderItemsToWarehouses(ctx context.Context, order *Order) error {
	orderItemStocks, err := d.orderItemStockRepository.GetListByOrderID(ctx, order.ID)
	if err != nil {
		return err
	}

	for _, orderItemStock := range orderItemStocks {
		stock, err := d.stockRepository.GetByWarehouseIDAndSku(ctx, orderItemStock.WarehouseID, orderItemStock.Sku)
		if err != nil {
			return err
		}

		stock.Count = stock.Count + orderItemStock.Count
		err = d.stockRepository.UpdateCount(ctx, stock)
		if err != nil {
			return err
		}

		err = d.orderItemStockRepository.Delete(ctx, orderItemStock)
		if err != nil {
			return err
		}
	}

	return nil
}
