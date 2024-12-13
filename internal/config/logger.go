package config

import (
	"os"
)

type LoggerConfig struct {
	Format    string `mapstructure:"Format"`
	Level     string `mapstructure:"Level"`
	Directory string `mapstructure:"Directory"`
	Name      string `mapstructure:"Name"`
}

func LoadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Format:    os.Getenv("LOGS_FORMAT"),
		Level:     os.Getenv("LOGS_LEVEL"),
		Directory: os.Getenv("LOGS_DIRECTORY"),
		Name:      os.Getenv("LOGS_NAME"),
	}
}
