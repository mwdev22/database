package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cli *redis.Client
}

func NewRedisCache(opt *redis.Options) *Redis {
	cli := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	return &Redis{
		cli: cli,
	}
}

func (c *Redis) Get(ctx context.Context, key string) (string, error) {
	return c.cli.Get(ctx, key).Result()
}

func (c *Redis) GetBytes(ctx context.Context, key string) ([]byte, error) {
	value, err := c.cli.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(value), nil
}

func (c *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.cli.Set(ctx, key, value, expiration).Err()
}

func (c *Redis) SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return c.cli.Set(ctx, key, value, expiration).Err()
}

func (c *Redis) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.cli.Keys(ctx, pattern).Result()
}

func (c *Redis) Del(ctx context.Context, key string) error {
	return c.cli.Del(ctx, key).Err()
}

func (c *Redis) Exists(ctx context.Context, key string) bool {
	return c.cli.Exists(ctx, key).Val() > 0
}
