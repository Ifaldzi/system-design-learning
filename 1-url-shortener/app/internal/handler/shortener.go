package handler

import (
	"net/http"

	"urlshortener/internal/service"
)

// Shortener exposes the URL shortener HTTP endpoints.
type Shortener struct {
	svc service.Shortener
}

func NewShortener(svc service.Shortener) *Shortener {
	return &Shortener{svc: svc}
}

// Create handles POST /api/v1/shorten.
// Expects: {"url": "...", "custom_alias": "optional"}.
// Returns: 201 with {"short_url", "short_key"} | 400 invalid input | 409 alias taken.
func (h *Shortener) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: decode body, call h.svc.Shorten, map errors, return 201 JSON.
	notImplemented(w)
}

// Redirect handles GET /{short_key}.
// Returns: 302/307 with Location header | 404 if unknown or expired.
func (h *Shortener) Redirect(w http.ResponseWriter, r *http.Request) {
	// TODO: read chi.URLParam(r, "short_key"), call h.svc.Resolve, redirect or 404.
	notImplemented(w)
}

// Analytics handles GET /api/v1/analytics/{short_key}.
// Returns: 200 with link metadata | 404 if unknown.
func (h *Shortener) Analytics(w http.ResponseWriter, r *http.Request) {
	// TODO: read chi.URLParam(r, "short_key"), call h.svc.Analytics, return 200 JSON.
	notImplemented(w)
}

func notImplemented(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_, _ = w.Write([]byte(`{"error":"not implemented"}`))
}
