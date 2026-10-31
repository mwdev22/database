package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const scanCount = 100

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
	// SCAN can return a key more than once, KEYS never does
	seen := make(map[string]struct{})
	keys := []string{}
	iter := c.cli.Scan(ctx, 0, pattern, scanCount).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		keys = append(keys, key)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.cli.Del(ctx, key).Err()
}

// fails open: a redis error reads as a miss, use Has where that matters
func (c *RedisCache) Exists(ctx context.Context, key string) bool {
	ok, _ := c.Has(ctx, key)
	return ok
}

func (c *RedisCache) Has(ctx context.Context, key string) (bool, error) {
	n, err := c.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
