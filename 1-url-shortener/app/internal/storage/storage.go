package storage

import (
	"context"
	"errors"
	"sync"
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

// Memory is an in-memory Storage implementation. The method bodies are left as
// stubs for you to implement.
type Memory struct {
	mu    sync.RWMutex
	links map[string]*Link
}

func NewMemory() *Memory {
	return &Memory{links: make(map[string]*Link)}
}

// Save persists a new link. It should return ErrConflict when the ShortKey
// already exists.
func (m *Memory) Save(ctx context.Context, link *Link) error {
	// TODO: implement
	return errors.New("storage.Memory.Save: not implemented")
}

// Get returns the link for shortKey or ErrNotFound. Expired links must not be
// returned.
func (m *Memory) Get(ctx context.Context, shortKey string) (*Link, error) {
	// TODO: implement
	return nil, ErrNotFound
}

// IncrementClicks atomically increments the click counter for shortKey.
func (m *Memory) IncrementClicks(ctx context.Context, shortKey string) error {
	// TODO: implement
	return errors.New("storage.Memory.IncrementClicks: not implemented")
}
