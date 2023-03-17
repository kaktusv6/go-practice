package domain

import "context"

func (d *domain) GetStocksBySKU(ctx context.Context, sku uint32) ([]Stock, error) {
	return d.stockRepository.GetListBySKU(ctx, sku)
}
