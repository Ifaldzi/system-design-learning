package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// OpenSQLite connects to a SQLite database using the pure-Go modernc driver.
// Use ":memory:" for an ephemeral DB or a path like "file:urlshortener.db".
func OpenSQLite(ctx context.Context, dsn string) (*SQLStore, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: open sqlite: %w", err)
	}
	configurePool(db)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("storage: ping sqlite: %w", err)
	}
	return NewSQLStore(db, DialectSQLite), nil
}
