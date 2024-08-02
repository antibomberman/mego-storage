package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	JWTSecret                string `env:"JWT_SECRET" required:"true"`
	StorageServiceServerPort string `env:"STORAGE_SERVICE_SERVER_PORT" required:"true"`
	MinioEndpoint            string `env:"MINIO_ENDPOINT" required:"true"`
	MinioAccessKey           string `env:"MINIO_ACCESS_KEY" required:"true"`
	MinioSecret              string `env:"MINIO_SECRET_KEY" required:"true"`
	MinioBucket              string `env:"MINIO_BUCKET" required:"true"`
	MinioPort                string `env:"MINIO_PORT" required:"true"`
	MinioRootUser            string `env:"MINIO_ROOT_USER" required:"true"`
	MinioRootPassword        string `env:"MINIO_ROOT_PASSWORD" required:"true"`
}

func Load() *Config {
	path := "./.env"
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	return &cfg
}
