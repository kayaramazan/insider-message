package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/kayaramazan/insider-message/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRedisCache(cfg *config.RedisConfig) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client, expiration: time.Duration(cfg.Expiration) * time.Hour}, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	return r.client.Set(ctx, key, value, r.expiration).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result > 0, err
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
