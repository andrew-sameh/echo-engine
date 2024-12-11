package config

import "os"

type ServerConfig struct {
	Host string
	Port string
	Env  string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
		Env:  os.Getenv("ENV"),
	}
}
