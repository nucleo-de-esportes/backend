package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type DatabaseConfig struct {
	Driver   string `env:"DRIVER"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	Name     string `env:"NAME"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD"`
}

type ServerConfig struct {
	Port string `env:"PORT" envDefault:"8000"`
}

type Config struct {
	DB     DatabaseConfig `envPrefix:"DB_"`
	Server ServerConfig   `envPrefix:"SERVER_"`
}

func LoadConfig() *Config {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("erro ao carregar variáveis de ambiente: %v", err)
	}

	return &cfg
}
