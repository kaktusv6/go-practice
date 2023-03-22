package repositories

import "time"

type Stock struct {
	Sku         uint32 `db:"sku"`
	WarehouseID int64  `db:"warehouse_id"`
	Count       uint64 `db:"count"`
}

type Order struct {
	ID        int64      `db:"id"`
	Status    string     `db:"status"`
	User      int64      `db:"user_id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type OrderItem struct {
	OrderId int64  `db:"order_id"`
	Sku     uint32 `db:"sku"`
	Count   uint16 `db:"count"`
}

type OrderItemStock struct {
	OrderId     int64  `db:"order_id"`
	Sku         uint32 `db:"sku"`
	Count       uint64 `db:"count"`
	WarehouseID int64  `db:"warehouse_id"`
}
