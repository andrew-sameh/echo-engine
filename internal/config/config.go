package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth   AuthConfig
	DB     DBConfig
	Server ServerConfig
	Redis  RedisConfig
	Logger LoggerConfig
	// LogLevel string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Auth:   LoadAuthConfig(),
		DB:     LoadDBConfig(),
		Server: LoadServerConfig(),
		Redis:  LoadRedisConfig(),
		Logger: LoadLoggerConfig(),
		// LogLevel: os.Getenv("LOG_LEVEL"),
	}
}
