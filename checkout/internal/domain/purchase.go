package domain

import (
	"context"
)

import (
	"github.com/pkg/errors"
	lomsV1Client "route256/loms/pkg/loms_v1"
)

func (d *domain) Purchase(ctx context.Context, user int64) error {
	userCartItems, err := d.cartItemRepository.GetUserCartItems(ctx, user)
	if err != nil {
		return err
	}

	if len(userCartItems) == 0 {
		return errors.New("user cart is empty")
	}

	itemInfoList := make([]*lomsV1Client.ItemInfo, 0, len(userCartItems))
	for _, cartItem := range userCartItems {
		itemInfoList = append(itemInfoList, &lomsV1Client.ItemInfo{
			Sku:   cartItem.Sku,
			Count: int32(cartItem.Count),
		})
	}

	_, err = d.lomsClient.CreateOrder(ctx, &lomsV1Client.OrderDataRequest{
		User:  user,
		Items: itemInfoList,
	})

	return err
}
