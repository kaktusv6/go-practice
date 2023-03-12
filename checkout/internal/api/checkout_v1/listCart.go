package checkoutV1

import (
	"context"
)

import (
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) ListCart(ctx context.Context, userInfo *desc.UserInfo) (*desc.CartInfo, error) {
	cart, err := i.domain.GetListItems(ctx, userInfo.GetUser())
	if err != nil {
		return nil, err
	}

	items := make([]*desc.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, &desc.CartItem{
			Sku:   item.Sku,
			Count: int32(item.Count),
			Name:  item.Name,
			Price: item.Price,
		})
	}

	cartInfo := &desc.CartInfo{
		Items:      items,
		TotalPrice: cart.TotalPrice,
	}

	return cartInfo, nil
}
