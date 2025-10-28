package database

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	GetBytes(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
}
