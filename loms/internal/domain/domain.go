package domain

type OrderCreator interface {
	CreateOrder(user int64, items []Item) (int64, error)
}

type ListOrderGetter interface {
	GetListOrder(orderID int64) (Order, error)
}

type OrderPayedMarker interface {
	OrderPayedMark(orderID int64) error
}

type StocksGetter interface {
	GetStocksBySKU(sku uint32) ([]Stock, error)
}

type OrderCanceling interface {
	CancelOrder(orderID int64) error
}

const (
	New             = "new"
	AwaitingPayment = "awaiting payment"
	Failed          = "failed"
	Payed           = "payed"
	Cancelled       = "cancelled"
)

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

type Domain struct {
}

func NewDomain() *Domain {
	return &Domain{}
}
