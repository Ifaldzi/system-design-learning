package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"urlshortener/internal/storage"
)

// ErrMiss is returned by Get when a key is absent or expired.
var ErrMiss = errors.New("cache miss")

// Cache sits in front of storage to absorb hot reads and negative lookups.
// Implementations must be safe for concurrent use.
type Cache interface {
	Get(ctx context.Context, shortKey string) (*storage.Link, error)
	Set(ctx context.Context, shortKey string, link *storage.Link, ttl time.Duration) error
	Delete(ctx context.Context, shortKey string) error
}

type entry struct {
	link      *storage.Link
	expiresAt time.Time
}

// Memory is an in-memory Cache implementation. Method bodies are stubs.
type Memory struct {
	mu    sync.RWMutex
	items map[string]entry
}

func NewMemory() *Memory {
	return &Memory{items: make(map[string]entry)}
}

// Get returns the cached link or ErrMiss. Expired entries are a miss.
func (m *Memory) Get(ctx context.Context, shortKey string) (*storage.Link, error) {
	// TODO: implement
	return nil, ErrMiss
}

// Set stores link under shortKey for ttl. A ttl <= 0 may mean "no expiry".
func (m *Memory) Set(ctx context.Context, shortKey string, link *storage.Link, ttl time.Duration) error {
	// TODO: implement
	return errors.New("cache.Memory.Set: not implemented")
}

// Delete removes shortKey if present. Deleting a missing key is not an error.
func (m *Memory) Delete(ctx context.Context, shortKey string) error {
	// TODO: implement
	return errors.New("cache.Memory.Delete: not implemented")
}
