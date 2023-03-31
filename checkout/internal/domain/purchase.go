package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
)

var (
	CartItemsEmptyError = errors.New("user cart is empty")
)

func (d *domain) Purchase(ctx context.Context, user int64) error {
	userCartItems, err := d.cartItemRepository.GetUserCartItems(ctx, user)
	if err != nil {
		return err
	}

	if len(userCartItems) == 0 {
		return CartItemsEmptyError
	}

	err = d.orderRepository.Create(ctx, user, userCartItems)
	if err != nil {
		return err
	}

	return nil
}
