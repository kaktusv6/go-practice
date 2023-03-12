package productServiceClient

type getProductRequestBody struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type listSKUsRequestBody struct {
	Token         string `json:"token"`
	StartAfterSku uint32 `json:"startAfterSku"`
	Count         uint32 `json:"count"`
}
