package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// OpenMySQL connects to MySQL/MariaDB using the go-sql-driver DSN, e.g.
// "user:pass@tcp(localhost:3306)/urlshortener?parseTime=true".
func OpenMySQL(ctx context.Context, dsn string) (*SQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("storage: open mysql: %w", err)
	}
	configurePool(db)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("storage: ping mysql: %w", err)
	}
	return NewSQLStore(db, DialectMySQL), nil
}
