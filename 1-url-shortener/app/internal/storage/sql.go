package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Dialect identifies the SQL dialect, mainly for placeholder formatting.
type Dialect string

const (
	DialectPostgres Dialect = "postgres"
	DialectMySQL    Dialect = "mysql"
	DialectSQLite   Dialect = "sqlite"
)

// placeholder returns the bind parameter for position n (1-indexed).
func (d Dialect) placeholder(n int) string {
	if d == DialectPostgres {
		return fmt.Sprintf("$%d", n)
	}
	return "?"
}

// SQLStore is a database/sql backed Storage implementation shared by every
// SQL dialect. Only the CRUD method bodies are left for you to implement.
type SQLStore struct {
	db      *sql.DB
	dialect Dialect
}

func NewSQLStore(db *sql.DB, dialect Dialect) *SQLStore {
	return &SQLStore{db: db, dialect: dialect}
}

// DB exposes the underlying pool for ad-hoc queries and migrations.
func (s *SQLStore) DB() *sql.DB { return s.db }

func (s *SQLStore) Close() error { return s.db.Close() }

func (s *SQLStore) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

// Save inserts link and returns ErrConflict when ShortKey already exists.
func (s *SQLStore) Save(ctx context.Context, link *Link) error {
	// TODO: implement INSERT using s.dialect.placeholder(1..n)
	return errors.New("storage.SQLStore.Save: not implemented")
}

// Get loads a link by key. Expired links must return ErrNotFound.
func (s *SQLStore) Get(ctx context.Context, shortKey string) (*Link, error) {
	// TODO: implement SELECT ... WHERE short_key = s.dialect.placeholder(1)
	return nil, ErrNotFound
}

// IncrementClicks atomically increments click_count for shortKey.
func (s *SQLStore) IncrementClicks(ctx context.Context, shortKey string) error {
	// TODO: implement UPDATE ... SET click_count = click_count + 1
	return errors.New("storage.SQLStore.IncrementClicks: not implemented")
}

// configurePool applies conservative pool defaults shared by all dialects.
func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)
}
