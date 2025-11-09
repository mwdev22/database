package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	config "github.com/mwdev22/gocfg"
)

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

	conf.MaxConns = int32(cfg.MaxOpenConns)
	conf.MinConns = int32(cfg.MinIdleConns)
	conf.MaxConnLifetime = time.Duration(cfg.ConnMaxLifetime) * time.Minute

	dbpool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create db pool: %w", err)
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return dbpool, nil
}
