package lomsV1

import (
	"context"
)

import (
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) ListOrder(ctx context.Context, req *desc.OrderID) (*desc.OrderResponse, error) {
	order, err := i.domain.GetListOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, err
	}
	items := make([]*desc.ItemInfo, 0, len(order.Items))

	for _, item := range order.Items {
		items = append(items, &desc.ItemInfo{
			Sku:   item.Sku,
			Count: int32(item.Count),
		})
	}

	status := desc.Status(desc.Status_value[order.Status])

	return &desc.OrderResponse{
		Status: status,
		User:   order.User,
		Items:  items,
	}, nil
}
