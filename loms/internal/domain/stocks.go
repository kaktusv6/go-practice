package domain

import "context"

func (d *domain) GetStocksBySKU(ctx context.Context, sku uint32) ([]Stock, error) {
	// Fixture
	return []Stock{
		{
			WarehouseID: 1,
			Count:       2,
		},
		{
			WarehouseID: 2,
			Count:       3,
		},
		{
			WarehouseID: 3,
			Count:       1,
		},
	}, nil
}
