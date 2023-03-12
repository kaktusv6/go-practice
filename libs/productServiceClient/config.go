package productServiceClient

type Config interface {
	GetUrl() string
	GetToken() string
}

type ConfigImpl struct {
	url   string
	token string
}

func (c *ConfigImpl) GetUrl() string {
	return c.url
}

func (c *ConfigImpl) GetToken() string {
	return c.token
}

func NewConfig(url, token string) Config {
	return &ConfigImpl{url, token}
}
