package receivers

type OrderStatusNotification struct {
	OrderID int64  `json:"orderId"`
	Status  string `json:"Status"`
}
