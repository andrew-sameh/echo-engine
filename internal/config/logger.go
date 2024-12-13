package config

import (
	"os"
	"strconv"
)

type LoggerConfig struct {
	Format    string `mapstructure:"Format"`
	Level     string `mapstructure:"Level"`
	Directory string `mapstructure:"Directory"`
	Name      string `mapstructure:"Name"`
	Local     bool   `mapstructure:"Local"`
}

func LoadLoggerConfig() LoggerConfig {
	local, err := strconv.ParseBool(os.Getenv("LOGS_LOCAL"))
	if err != nil {
		local = false
	}

	return LoggerConfig{
		Format:    os.Getenv("LOGS_FORMAT"),
		Level:     os.Getenv("LOGS_LEVEL"),
		Directory: os.Getenv("LOGS_DIRECTORY"),
		Name:      os.Getenv("LOGS_NAME"),
		Local:     local,
	}
}
