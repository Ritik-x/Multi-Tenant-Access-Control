package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	AppEnv string `env:"APP_ENV" env-default:"development"`
	Port   string `env:"PORT" env-default:"8080"`

	DBHost     string `env:"DB_HOST" env-required:"true"`
	DBPort     string `env:"DB_PORT" env-required:"true"`
	DBUser     string `env:"DB_USER" env-required:"true"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME" env-required:"true"`
	DBSSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}
func Load() *Config{
var cfg Config
if err := cleanenv.ReadEnv(&cfg) ; err != nil {
			log.Fatal("failed to load configuration: ", err)
}
return &cfg
}