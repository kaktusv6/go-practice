package lomsClient

type StocksResBody struct {
	Stocks []Stock `json:"stocks"`
}

type createOrderResBody struct {
	OrderID int64 `json:"orderID"'`
}
