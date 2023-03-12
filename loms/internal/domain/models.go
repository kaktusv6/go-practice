package domain

type Item struct {
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Order struct {
	Status string `json:"status"`
	User   int64  `json:"user"`
	Items  []Item `json:"items"`
}

type Stock struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}
