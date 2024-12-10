package config

import "os"

type ServerConfig struct {
	Host       string
	Port       string
	ExposePort string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Host:       os.Getenv("HOST"),
		Port:       os.Getenv("PORT"),
		ExposePort: os.Getenv("EXPOSE_PORT"),
	}
}
