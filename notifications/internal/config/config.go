package config

type Config struct {
	App     AppConfig `yaml:"app"`
	Brokers []string  `yaml:"brokers"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}
