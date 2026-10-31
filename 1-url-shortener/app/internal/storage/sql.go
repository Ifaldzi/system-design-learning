package storage

import (
	"context"
	"database/sql"
	"time"
)

// DBStore shares the database pool and lifecycle methods across SQL backends.
// Backend-specific stores embed it and implement their own CRUD methods.
type DBStore struct {
	db *sql.DB
}

func newDBStore(ctx context.Context, driver, dsn string) (*DBStore, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	configurePool(db)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &DBStore{db: db}, nil
}

// DB exposes the underlying pool for ad-hoc queries and migrations.
func (s *DBStore) DB() *sql.DB { return s.db }

func (s *DBStore) Close() error { return s.db.Close() }

func (s *DBStore) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

// configurePool applies conservative pool defaults shared by all dialects.
func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
}
