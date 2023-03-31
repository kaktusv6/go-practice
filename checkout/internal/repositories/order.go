package repositories

import (
	"context"
)

import (
	"route256/checkout/internal/domain"
	lomsV1Clinet "route256/loms/pkg/loms_v1"
)

type OrderRepository struct {
	lomsClient lomsV1Clinet.LomsV1Client
}

func NewOrderRepository(lomsClient lomsV1Clinet.LomsV1Client) domain.OrderRepository {
	return &OrderRepository{
		lomsClient: lomsClient,
	}
}

func (o *OrderRepository) Create(ctx context.Context, user int64, items []*domain.CartItem) error {
	itemInfoList := make([]*lomsV1Clinet.ItemInfo, 0, len(items))
	for _, cartItem := range items {
		itemInfoList = append(itemInfoList, &lomsV1Clinet.ItemInfo{
			Sku:   cartItem.Sku,
			Count: int32(cartItem.Count),
		})
	}

	_, err := o.lomsClient.CreateOrder(ctx, &lomsV1Clinet.OrderDataRequest{
		User:  user,
		Items: itemInfoList,
	})

	if err != nil {
		return err
	}

	return nil
}
