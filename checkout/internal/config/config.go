package config

type Config struct {
	App            AppConfig      `yaml:"app"`
	ProductService ProductService `yaml:"product_service"`
	Loms           Loms           `yaml:"loms"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type ProductService struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type Loms struct {
	Url string `yaml:"url"`
}
