package config

type Config struct {
	App     AppConfig `yaml:"app"`
	Brokers []string  `yaml:"brokers"`
	Logger  Logger    `yaml:"logger"`
}

type AppConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type Logger struct {
	Level string `yaml:"level"`
}
