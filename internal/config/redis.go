package config

import (
	"fmt"
	"os"
	"strconv"
)

type RedisConfig struct {
	Host      string `mapstructure:"Host"`
	Port      int    `mapstructure:"Port"`
	Password  string `mapstructure:"Password"`
	KeyPrefix string `mapstructure:"KeyPrefix"`
	MainDB    int    `mapstructure:"MainDB"`
	TaskDB    int    `mapstructure:"TaskDB"`
}

func LoadRedisConfig() RedisConfig {
	port, _ := strconv.Atoi(os.Getenv("REDIS_PORT"))
	mainDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	taskDB, _ := strconv.Atoi(os.Getenv("REDIS_TASK_DB"))
	return RedisConfig{
		Host:      os.Getenv("REDIS_HOST"),
		Port:      port,
		Password:  os.Getenv("REDIS_PASSWORD"),
		KeyPrefix: os.Getenv("REDIS_KEY_PREFIX"),
		MainDB:    mainDB,
		TaskDB:    taskDB,
	}
}
func (a *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
