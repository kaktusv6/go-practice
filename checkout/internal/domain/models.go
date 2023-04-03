package domain

type Cart struct {
	Items      []*CartItem
	TotalPrice uint64
}

func (c *Cart) CalculateTotalPrice() {
	var result uint64 = 0
	for _, item := range c.Items {
		if item.Product != nil {
			result += uint64(item.Product.Price) * uint64(item.Count)
		}
	}
	c.TotalPrice = result
}

type CartItem struct {
	User    int64
	Sku     uint32
	Count   uint16
	Product *ProductInfo
}

type ProductInfo struct {
	Name  string
	Price uint32
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}
