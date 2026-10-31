package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	config "github.com/mwdev22/gocfg"
)

const connectTimeout = 10 * time.Second

type Pgx struct {
	Config *config.DatabaseConfig
	Pool   *pgxpool.Pool
}

func New(cfg *config.DatabaseConfig) (*Pgx, error) {
	return &Pgx{
		Config: cfg,
	}, nil
}

func (pg *Pgx) Connect() error {
	pool, err := connect(pg.Config)
	if err != nil {
		return err
	}
	pg.Pool = pool
	return nil
}

func connect(cfg *config.DatabaseConfig) (*pgxpool.Pool, error) {
	configStr := cfg.URI

	conf, err := pgxpool.ParseConfig(configStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse db config: %w", err)
	}

	applyPoolConfig(conf, cfg)

	// the pool keeps this ctx for creating min conns in the background, so it must outlive connect
	dbpool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return dbpool, nil
}

// zero values keep pgx defaults and pool_* params from the uri
func applyPoolConfig(conf *pgxpool.Config, cfg *config.DatabaseConfig) {
	if cfg.MaxOpenConns > 0 {
		conf.MaxConns = int32(cfg.MaxOpenConns)
	}
	if cfg.MinIdleConns > 0 {
		conf.MinConns = int32(cfg.MinIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		conf.MaxConnLifetime = time.Duration(cfg.ConnMaxLifetime) * time.Minute
	}
}
