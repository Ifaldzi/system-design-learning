package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// OpenPostgres connects to PostgreSQL using a libpq DSN, e.g.
// "postgres://user:pass@localhost:5432/urlshortener?sslmode=disable".
func OpenPostgres(ctx context.Context, dsn string) (*SQLStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: open postgres: %w", err)
	}
	configurePool(db)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("storage: ping postgres: %w", err)
	}
	return NewSQLStore(db, DialectPostgres), nil
}
