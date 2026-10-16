package service

import (
	"context"
	"errors"

	"urlshortener/internal/cache"
	"urlshortener/internal/storage"
)

var (
	ErrInvalidURL   = errors.New("invalid url")
	ErrInvalidAlias = errors.New("invalid alias")
	ErrAliasTaken   = errors.New("alias already taken")
	ErrNotFound     = storage.ErrNotFound
)

// ShortenRequest is the decoded body of POST /api/v1/shorten.
type ShortenRequest struct {
	URL         string `json:"url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}

// ShortenResult is returned on a successful shorten.
type ShortenResult struct {
	ShortURL string `json:"short_url"`
	ShortKey string `json:"short_key"`
}

// Shortener is the core business-logic contract used by the HTTP layer.
type Shortener interface {
	Shorten(ctx context.Context, req ShortenRequest) (*ShortenResult, error)
	Resolve(ctx context.Context, shortKey string) (*storage.Link, error)
	Analytics(ctx context.Context, shortKey string) (*storage.Link, error)
}

type shortener struct {
	store   storage.Storage
	cache   cache.Cache
	baseURL string
}

func NewShortener(store storage.Storage, c cache.Cache, baseURL string) Shortener {
	return &shortener{store: store, cache: c, baseURL: baseURL}
}

// Shorten validates the input, generates/accepts a key and persists the link.
func (s *shortener) Shorten(ctx context.Context, req ShortenRequest) (*ShortenResult, error) {
	// TODO: implement
	return nil, errors.New("service.Shorten: not implemented")
}

// Resolve looks up a key, increments its click counter and returns the link.
func (s *shortener) Resolve(ctx context.Context, shortKey string) (*storage.Link, error) {
	// TODO: implement
	return nil, ErrNotFound
}

// Analytics returns link metadata for the given key.
func (s *shortener) Analytics(ctx context.Context, shortKey string) (*storage.Link, error) {
	// TODO: implement
	return nil, ErrNotFound
}
