package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const ConfigPath = "./config/config.json"

type Config struct {
	HTTP HTTPConfig `json:"http" env-prefix:"HTTP_"`
	DB   DBConfig   `json:"db" env-prefix:"DB_"`
	App  AppConfig  `json:"app" env-prefix:"APP_"`
}

type HTTPConfig struct {
	Port              string        `json:"port" env:"PORT" env-default:"8081"`
	Address           string        `json:"address" env:"ADDRESS" env-default:"localhost"`
	ReadTimeout       time.Duration `json:"read_timeout" env:"READ_TIMEOUT" env-default:"10s"`
	WriteTimeout      time.Duration `json:"write_timeout" env:"WRITE_TIMEOUT" env-default:"10s"`
	ReadHeaderTimeout time.Duration `json:"read_header_timeout" env:"READ_HEADER_TIMEOUT" env-default:"5s"`
	IdleTimeout       time.Duration `json:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

type DBConfig struct {
	Host     string `json:"host" env:"HOST" env-default:"localhost"`
	Port     string `json:"port" env:"PORT" env-default:"5434"`
	User     string `json:"user" env:"USER" env-default:"payment_gateway"`
	Password string `json:"password" env:"PASSWORD"`
	Name     string `json:"name" env:"NAME" env-default:"payment_gateway"`
	SSLMode  string `json:"ssl_mode" env:"SSL_MODE" env-default:"disable"`
}

type AppConfig struct {
	Env string `json:"env" env:"ENV" env-default:"prod"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
