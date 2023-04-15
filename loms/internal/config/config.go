package config

type Config struct {
	App            AppConfig      `yaml:"app"`
	ProductService ProductService `yaml:"product_service"`
	DataBase       DataBase       `yaml:"database"`
	Brokers        []string       `yaml:"brokers"`
	Logger         Logger         `yaml:"logger"`
}

type AppConfig struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"env"`
	Port        string `yaml:"port"`
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

type Logger struct {
	Level string `yml:"level"`
}
