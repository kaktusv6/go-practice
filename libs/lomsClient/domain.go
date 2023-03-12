package lomsClient

type Stock struct {
	WarehouseID int64  `json:"warehouseID"`
	Count       uint64 `json:"count"`
}

type Item struct {
	Sku   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}
