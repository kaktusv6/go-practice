package config

type Config struct {
	App            AppConfig      `yaml:"app"`
	ProductService ProductService `yaml:"product_service"`
	DataBase       DataBase       `yaml:"database"`
	Brokers        []string       `yaml:"brokers"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type ProductService struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}
