package storage

import (
	"context"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

// OpenSQLite connects to a SQLite database using the pure-Go modernc driver.
// Use ":memory:" for an ephemeral DB or a path like "file:urlshortener.db".
func OpenSQLite(ctx context.Context, dsn string) (*SQLiteStore, error) {
	db, err := newDBStore(ctx, "sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: open sqlite: %w", err)
	}
	return &SQLiteStore{DBStore: db}, nil
}

type SQLiteStore struct {
	*DBStore
}

// Save inserts a link. A duplicate key is reported as ErrConflict.
func (s *SQLiteStore) Save(ctx context.Context, link *Link) error {
	// TODO: implement the SQLite INSERT using ? placeholders.
	return errors.New("storage.SQLiteStore.Save: not implemented")
}

// Get loads a non-expired link by key.
func (s *SQLiteStore) Get(ctx context.Context, shortKey string) (*Link, error) {
	// TODO: implement the SQLite SELECT using ? for shortKey.
	return nil, errors.New("storage.SQLiteStore.Get: not implemented")
}

// IncrementClicks atomically increments click_count for shortKey.
func (s *SQLiteStore) IncrementClicks(ctx context.Context, shortKey string) error {
	// TODO: implement an atomic SQLite UPDATE using ? for shortKey.
	return errors.New("storage.SQLiteStore.IncrementClicks: not implemented")
}
