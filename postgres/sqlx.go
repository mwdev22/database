package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	config "github.com/mwdev22/gocfg"
)

type Sqlx struct {
	Cfg *config.DatabaseConfig
	Db  *sqlx.DB
}

func NewSqlx(cfg *config.DatabaseConfig) (*Pgx, error) {
	return &Pgx{
		Config: cfg,
	}, nil
}

func (s *Sqlx) Connect() error {
	db, err := s.connect(s.Cfg)
	if err != nil {
		return err
	}
	s.Db = db
	return nil
}

func (s *Sqlx) connect(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlx db: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	// test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping sqlx db: %w", err)
	}

	return db, nil
}
