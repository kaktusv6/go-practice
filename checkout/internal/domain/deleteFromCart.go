package domain

import "context"

func (d *domain) DeleteFromCart(ctx context.Context, cartItem CartItem) error {
	existCartItem, err := d.cartItemRepository.GetOne(ctx, cartItem.User, cartItem.Sku)
	if err != nil {
		return err
	}

	existCartItem.Count = existCartItem.Count - cartItem.Count
	err = d.transactionManager.RunRepeatableReade(ctx, func(ctxTx context.Context) error {
		if existCartItem.Count == 0 {
			err = d.cartItemRepository.Delete(ctxTx, existCartItem)
		} else {
			err = d.cartItemRepository.Update(ctxTx, existCartItem)
		}
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
