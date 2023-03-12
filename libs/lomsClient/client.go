package lomsClient

import "route256/libs/httpClientWrapper"

type Client interface {
	Stocks(sku uint32) ([]Stock, error)
	CreateOrder(user int64, items []Item) (int64, error)
}

type ClientImpl struct {
	config Config
}

func New(config Config) Client {
	return &ClientImpl{config}
}

func (c *ClientImpl) Stocks(sku uint32) ([]Stock, error) {
	req := stocksReqBody{
		sku,
	}

	res, err := httpClientWrapper.New[stocksReqBody, StocksResBody](
		c.config.GetUrl() + "/stocks",
	).SendRequest(req)

	return res.Stocks, err
}

func (c *ClientImpl) CreateOrder(user int64, items []Item) (int64, error) {
	req := createOrderReqBody{
		user,
		items,
	}

	res, err := httpClientWrapper.New[createOrderReqBody, createOrderResBody](
		c.config.GetUrl() + "/createOrder",
	).SendRequest(req)

	return res.OrderID, err
}
