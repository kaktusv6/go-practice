package lomsClient

type stocksReqBody struct {
	SKU uint32 `json:"sku"`
}

type createOrderReqBody struct {
	User  int64  `json:"user"`
	Items []Item `json:"items"`
}
