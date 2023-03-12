package lomsClient

type Config interface {
	GetUrl() string
}

type ConfigImpl struct {
	url string
}

func (c *ConfigImpl) GetUrl() string {
	return c.url
}

func NewConfig(url string) Config {
	return &ConfigImpl{url}
}
