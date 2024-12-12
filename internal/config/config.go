package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth     AuthConfig
	DB       DBConfig
	Server   ServerConfig
	Redis    RedisConfig
	LogLevel string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Auth:     LoadAuthConfig(),
		DB:       LoadDBConfig(),
		Server:   LoadServerConfig(),
		Redis:    LoadRedisConfig(),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
