package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/andrew-sameh/echo-engine/internal/config"
	"github.com/andrew-sameh/echo-engine/pkg/logger"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cache  *cache.Cache
	client *redis.Client
	prefix string
}

// NewRedis creates a new redis client instance
func NewRedis(cfg *config.Config) Redis {
	addr := cfg.Redis.Addr()

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       cfg.Redis.MainDB,
		Password: cfg.Redis.Password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Log().Fatalf("Error to open redis[%s] connection: %v", addr, err)
	}

	logger.Log().Info("Redis connection established")
	return Redis{
		client: client,
		prefix: cfg.Redis.KeyPrefix,
		cache: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		}),
	}
}

func (a Redis) wrapperKey(key string) string {
	return fmt.Sprintf("%s:%s", a.prefix, key)
}

func (a Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return a.cache.Set(&cache.Item{
		Ctx:            context.TODO(),
		Key:            a.wrapperKey(key),
		Value:          value,
		TTL:            expiration,
		SkipLocalCache: true,
	})
}

func (a Redis) Get(key string, value interface{}) error {
	err := a.cache.Get(context.TODO(), a.wrapperKey(key), value)
	if err == cache.ErrCacheMiss {
		err = fmt.Errorf("Key %s not found", key)
	}

	return err
}

func (a Redis) Delete(keys ...string) (bool, error) {
	wrapperKeys := make([]string, len(keys))
	for index, key := range keys {
		wrapperKeys[index] = a.wrapperKey(key)
	}

	cmd := a.client.Del(context.TODO(), wrapperKeys...)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() > 0, nil
}

func (a Redis) Check(keys ...string) (bool, error) {
	wrapperKeys := make([]string, len(keys))
	for index, key := range keys {
		wrapperKeys[index] = a.wrapperKey(key)
	}

	cmd := a.client.Exists(context.TODO(), wrapperKeys...)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val() > 0, nil
}

func (a Redis) Close() error {
	return a.client.Close()
}

func (a Redis) GetClient() *redis.Client {
	return a.client
}
