package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	GrpcPort string `env:"GRPC_PORT" env-default:":50051"`

	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" env-default:"10"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" env-default:"30"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" env-default:"10"`

	DbConfig
}

type DbConfig struct {
	DBHost     string `env:"DB_HOST" env-default:"localhost"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBUser     string `env:"DB_USER" env-required:"true"`
	DBPassword string `env:"DB_PASSWORD" env-required:"true"`
	DBDbname   string `env:"DB_NAME" env-required:"true"`
	DBSslmode  string `env:"DB_SSLMODE" env-default:"disable"`
}

func MustConfig() *Config {
	if err := godotenv.Load("services/orders/.env"); err != nil {
		log.Fatalf("load orders env: %v", err)
	}

	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("read orders config: %v", err)
	}

	return &cfg
}
