package repositories

type CartItem struct {
	User  int64  `db:"user_id"`
	Sku   uint32 `db:"sku"`
	Count uint16 `db:"count"`
}
