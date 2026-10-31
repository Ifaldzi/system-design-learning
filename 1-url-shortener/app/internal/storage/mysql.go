package storage

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// OpenMySQL connects to MySQL/MariaDB using the go-sql-driver DSN, e.g.
// "user:pass@tcp(localhost:3306)/urlshortener?parseTime=true".
func OpenMySQL(ctx context.Context, dsn string) (*MySQLStore, error) {
	db, err := newDBStore(ctx, "mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: open mysql: %w", err)
	}
	return &MySQLStore{DBStore: db}, nil
}

type MySQLStore struct {
	*DBStore
}

// Save inserts a link. A duplicate key is reported as ErrConflict.
func (s *MySQLStore) Save(ctx context.Context, link *Link) error {
	// TODO: implement the MySQL INSERT using ? placeholders.
	return errors.New("storage.MySQLStore.Save: not implemented")
}

// Get loads a non-expired link by key.
func (s *MySQLStore) Get(ctx context.Context, shortKey string) (*Link, error) {
	// TODO: implement the MySQL SELECT using ? for shortKey.
	return nil, errors.New("storage.MySQLStore.Get: not implemented")
}

// IncrementClicks atomically increments click_count for shortKey.
func (s *MySQLStore) IncrementClicks(ctx context.Context, shortKey string) error {
	// TODO: implement an atomic MySQL UPDATE using ? for shortKey.
	return errors.New("storage.MySQLStore.IncrementClicks: not implemented")
}
