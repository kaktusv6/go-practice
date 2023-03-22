package domain

import "time"

type Item struct {
	Sku   uint32
	Count uint16
}

type Order struct {
	ID        int64
	Status    string
	User      int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Items     []Item
}

type Stock struct {
	Sku         uint32
	WarehouseID int64
	Count       uint64
}

type OrderItemStock struct {
	OrderId     int64
	Sku         uint32
	Count       uint64
	WarehouseID int64
}
