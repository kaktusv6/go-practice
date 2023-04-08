package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
)

var (
	ErrorEmptyItems = errors.New("items is empty")
)

func (d *domain) CreateOrder(ctx context.Context, order *Order) (int64, error) {
	if len(order.Items) == 0 {
		return 0, ErrorEmptyItems
	}

	order.Status = New

	err := d.manager.RepeatableRead(ctx, func(ctxTx context.Context) error {
		// Сохраняем заказ и его товары
		err := d.saveOrder(ctxTx, order)
		if err != nil {
			return err
		}

		// Записываем резервирование товаров
		isSuccessReserve := true
		for _, item := range order.Items {
			// Выполняем резервацию каждой позиции в заказе
			if isSuccessReserve {
				err = d.reserveOrderItem(ctxTx, order, item)
			}

			// Если получили недостаток товаров то помечаем заказ как провальный
			if err != nil && errors.Is(err, ErrorNotEnoughItems) {
				order.Status = Failed
				isSuccessReserve = false
			} else if err != nil {
				return err
			}
		}

		// Если все хорошо то помечаем заказ как ожидающий оплаты
		if isSuccessReserve {
			order.Status = AwaitingPayment
		}

		// Обновляем данные о заказе
		err = d.orderRepository.Update(ctxTx, order)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	err = d.orderStatusNotifier.Notify(order)
	if err != nil {
		return 0, err
	}

	return order.ID, nil
}

func (d *domain) saveOrder(ctx context.Context, order *Order) error {
	// Сохраняем информацию о заказе
	err := d.orderRepository.Save(ctx, order)
	if err != nil {
		return err
	}

	// Сохраняем товары заказа
	err = d.orderItemRepository.SaveMany(ctx, order.ID, order.Items)
	if err != nil {
		return err
	}

	return nil
}

// Метод резервирует товар заказа
func (d *domain) reserveOrderItem(ctx context.Context, order *Order, item *Item) error {
	// Получаем остатки товара на всех складах
	stocks, err := d.stockRepository.GetListBySKU(ctx, item.Sku)
	if err != nil {
		return err
	}

	// Кол-во товара которое надо разервировать
	count := uint64(item.Count)
	// Проходимся по всем складам где есть остатки товара
	for _, stock := range stocks {
		newCount := count
		// Проверяем мы списываем все остатки со склада или только часть
		if count >= stock.Count {
			newCount = count - stock.Count
		} else {
			newCount = 0
		}

		// Формируем резервирование только если есть что резервировать
		if count > 0 {
			orderItemStockCount := stock.Count
			if newCount <= 0 {
				orderItemStockCount = count
			}

			err = d.orderItemStockRepository.Save(ctx, &OrderItemStock{
				OrderId:     order.ID,
				Sku:         item.Sku,
				Count:       orderItemStockCount,
				WarehouseID: stock.WarehouseID,
			})

			if err != nil {
				return err
			}
		}

		count = newCount
	}

	if count > 0 {
		return ErrorNotEnoughItems
	}

	return nil
}
