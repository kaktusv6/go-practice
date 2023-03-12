package domain

import (
	"context"
)

import (
	productServiceV1Clinet "route256/checkout/pkg/product_service_v1"
)

type Cart struct {
	Items      []Item `json:"items"`
	TotalPrice uint32 `json:"totalPrice"`
}

func (c *Cart) calculateTotalPrice() {
	var result uint32 = 0
	for _, item := range c.Items {
		result += item.Price
	}
	c.TotalPrice = result
}

type Item struct {
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (d *domain) GetListItems(ctx context.Context, user int64) (Cart, error) {
	skuItems := []uint32{
		1076963,
		1148162,
		1625903,
	}

	skuCounts := map[uint32]uint16{
		1076963: 1,
		1148162: 5,
		1625903: 2,
	}

	cart := Cart{}
	products := map[uint32]*productServiceV1Clinet.GetProductResponse{}
	for _, sku := range skuItems {
		product, err := d.productServiceClient.GetProduct(ctx, &productServiceV1Clinet.GetProductRequest{
			Token: d.productServiceToken,
			Sku:   sku,
		})
		if err != nil {
			return cart, err
		}
		products[sku] = product
	}

	items := make([]Item, 0, len(products))
	for sku, product := range products {
		item := Item{
			Sku:   sku,
			Count: skuCounts[sku],
			Name:  product.GetName(),
			Price: product.GetPrice(),
		}
		items = append(items, item)
	}
	cart.Items = items
	cart.calculateTotalPrice()

	return cart, nil
}
