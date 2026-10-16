package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("short key not found")
	ErrConflict = errors.New("short key already exists")
)

// Link is the persisted representation of a shortened URL.
type Link struct {
	ShortKey    string
	OriginalURL string
	ClickCount  int64
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

// Storage abstracts persistence so the service layer is agnostic to the backend.
type Storage interface {
	Save(ctx context.Context, link *Link) error
	Get(ctx context.Context, shortKey string) (*Link, error)
	IncrementClicks(ctx context.Context, shortKey string) error
}
