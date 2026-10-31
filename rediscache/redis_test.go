package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// no live server: a closed port makes every command fail, which is the case the fixes are about
func newUnreachable(t *testing.T) *RedisCache {
	t.Helper()
	c := New(&redis.Options{
		Addr:        "127.0.0.1:1",
		DialTimeout: 200 * time.Millisecond,
		MaxRetries:  -1,
	})
	t.Cleanup(func() { c.cli.Close() })
	return c
}

func TestRedisErrorsSurface(t *testing.T) {
	c := newUnreachable(t)
	ctx := context.Background()

	tests := []struct {
		name string
		run  func() error
	}{
		{name: "Has returns error", run: func() error {
			ok, err := c.Has(ctx, "k")
			if ok {
				t.Error("Has = true on error")
			}
			return err
		}},
		{name: "Keys returns error", run: func() error {
			keys, err := c.Keys(ctx, "*")
			if keys != nil {
				t.Errorf("Keys = %v on error, want nil", keys)
			}
			return err
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.run(); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestExistsFalseOnError(t *testing.T) {
	c := newUnreachable(t)
	if c.Exists(context.Background(), "k") {
		t.Error("Exists = true on error")
	}
}
