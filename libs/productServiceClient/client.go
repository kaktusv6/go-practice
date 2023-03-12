package productServiceClient

import (
	"route256/libs/httpClientWrapper"
)

type Client interface {
	GetProduct(sku uint32) (Product, error)
	ListSkus(startAfterSku uint32, count uint32) ([]uint32, error)
}

type ClientImpl struct {
	config Config
}

func New(config Config) Client {
	return &ClientImpl{config}
}

func (c *ClientImpl) GetProduct(sku uint32) (Product, error) {
	req := getProductRequestBody{
		c.config.GetToken(),
		sku,
	}

	res, err := httpClientWrapper.New[getProductRequestBody, Product](
		c.config.GetUrl() + "/get_product",
	).SendRequest(req)

	return *res, err
}

func (c *ClientImpl) ListSkus(startAfterSku uint32, count uint32) ([]uint32, error) {
	req := listSKUsRequestBody{
		c.config.GetToken(),
		startAfterSku,
		count,
	}

	res, err := httpClientWrapper.New[listSKUsRequestBody, listSKUSResponseBody](
		c.config.GetUrl() + "/get_product",
	).SendRequest(req)
	if err != nil {
		return res.SKUs, nil
	}

	return res.SKUs, nil
}
