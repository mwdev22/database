package postgres

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	config "github.com/mwdev22/gocfg"
)

func TestApplyPoolConfig(t *testing.T) {
	const uri = "postgres://u:p@localhost:5432/db?pool_max_conns=7&pool_min_conns=2&pool_max_conn_lifetime=30m"

	tests := []struct {
		name         string
		cfg          config.DatabaseConfig
		wantMax      int32
		wantMin      int32
		wantLifetime time.Duration
	}{
		{name: "zero values keep uri settings", cfg: config.DatabaseConfig{}, wantMax: 7, wantMin: 2, wantLifetime: 30 * time.Minute},
		{name: "positive values override", cfg: config.DatabaseConfig{MaxOpenConns: 20, MinIdleConns: 5, ConnMaxLifetime: 10}, wantMax: 20, wantMin: 5, wantLifetime: 10 * time.Minute},
		{name: "partial override", cfg: config.DatabaseConfig{MaxOpenConns: 12}, wantMax: 12, wantMin: 2, wantLifetime: 30 * time.Minute},
		{name: "negative values ignored", cfg: config.DatabaseConfig{MaxOpenConns: -1, MinIdleConns: -1, ConnMaxLifetime: -1}, wantMax: 7, wantMin: 2, wantLifetime: 30 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, err := pgxpool.ParseConfig(uri)
			if err != nil {
				t.Fatalf("parse config: %v", err)
			}
			cfg := tt.cfg
			cfg.URI = uri

			applyPoolConfig(conf, &cfg)

			if conf.MaxConns != tt.wantMax {
				t.Errorf("MaxConns = %d, want %d", conf.MaxConns, tt.wantMax)
			}
			if conf.MinConns != tt.wantMin {
				t.Errorf("MinConns = %d, want %d", conf.MinConns, tt.wantMin)
			}
			if conf.MaxConnLifetime != tt.wantLifetime {
				t.Errorf("MaxConnLifetime = %v, want %v", conf.MaxConnLifetime, tt.wantLifetime)
			}
		})
	}
}

func TestApplyPoolConfigKeepsPgxDefaults(t *testing.T) {
	conf, err := pgxpool.ParseConfig("postgres://u:p@localhost:5432/db")
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}
	want := conf.MaxConns

	applyPoolConfig(conf, &config.DatabaseConfig{})

	if conf.MaxConns != want || conf.MaxConns < 1 {
		t.Errorf("MaxConns = %d, want pgx default %d", conf.MaxConns, want)
	}
}
