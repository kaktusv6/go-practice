package lomsV1

import (
	"context"
)

import (
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) CreateOrder(ctx context.Context, req *desc.OrderDataRequest) (*desc.OrderID, error) {
	items := make([]domain.Item, 0, len(req.GetItems()))
	for _, reqItem := range req.GetItems() {
		items = append(items, domain.Item{
			Sku:   reqItem.GetSku(),
			Count: uint16(reqItem.GetCount()),
		})
	}

	orderID, err := i.domain.CreateOrder(ctx, req.GetUser(), items)
	if err != nil {
		return nil, err
	}
	return &desc.OrderID{
		OrderId: orderID,
	}, nil
}
