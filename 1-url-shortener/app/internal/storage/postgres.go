package storage

import (
	"context"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	*DBStore
}

// OpenPostgres connects to PostgreSQL using a libpq DSN, e.g.
// "postgres://user:pass@localhost:5432/urlshortener?sslmode=disable".
func OpenPostgres(ctx context.Context, dsn string) (*PostgresStore, error) {
	db, err := newDBStore(ctx, "postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: open postgres: %w", err)
	}
	return &PostgresStore{DBStore: db}, nil
}

// Save inserts a link. A duplicate key is reported as ErrConflict.
func (s *PostgresStore) Save(ctx context.Context, link *Link) error {
	// TODO: implement the PostgreSQL INSERT using $1, $2, ... placeholders.
	return errors.New("storage.PostgresStore.Save: not implemented")
}

// Get loads a non-expired link by key.
func (s *PostgresStore) Get(ctx context.Context, shortKey string) (*Link, error) {
	// TODO: implement the PostgreSQL SELECT using $1 for shortKey.
	return nil, errors.New("storage.PostgresStore.Get: not implemented")
}

// IncrementClicks atomically increments click_count for shortKey.
func (s *PostgresStore) IncrementClicks(ctx context.Context, shortKey string) error {
	// TODO: implement an atomic PostgreSQL UPDATE using $1 for shortKey.
	return errors.New("storage.PostgresStore.IncrementClicks: not implemented")
}
