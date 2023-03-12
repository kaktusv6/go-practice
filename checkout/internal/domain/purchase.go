package domain

import (
	"context"
)

import (
	lomsV1Client "route256/loms/pkg/loms_v1"
)

func (d *domain) Purchase(ctx context.Context, user int64) error {
	// Fixture
	fixtureItems := []*lomsV1Client.ItemInfo{
		{Sku: 1076963, Count: 1},
		{Sku: 1148162, Count: 4},
		{Sku: 1625903, Count: 2},
	}

	_, err := d.lomsClient.CreateOrder(ctx, &lomsV1Client.OrderDataRequest{
		User:  user,
		Items: fixtureItems,
	})

	return err
}
