package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cli *redis.Client
}

func New(opt *redis.Options) *RedisCache {
	cli := redis.NewClient(opt)
	return &RedisCache{
		cli: cli,
	}
}

func (c *RedisCache) Connect(opt *redis.Options) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, err := c.cli.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %s", err)
	}
	return nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.cli.Get(ctx, key).Result()
}

func (c *RedisCache) GetBytes(ctx context.Context, key string) ([]byte, error) {
	value, err := c.cli.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(value), nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.cli.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return c.cli.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.cli.Keys(ctx, pattern).Result()
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.cli.Del(ctx, key).Err()
}

func (c *RedisCache) Exists(ctx context.Context, key string) bool {
	return c.cli.Exists(ctx, key).Val() > 0
}
