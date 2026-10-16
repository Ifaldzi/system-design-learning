package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"urlshortener/internal/cache"
	"urlshortener/internal/handler"
	"urlshortener/internal/service"
	"urlshortener/internal/storage"
)

func main() {
	baseURL := env("BASE_URL", "http://localhost:8080")
	addr := env("ADDR", ":8080")

	store := storage.NewMemory()
	c := cache.NewMemory()
	svc := service.NewShortener(store, c, baseURL)
	h := handler.NewShortener(svc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Get("/health", handler.Health)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", h.Create)
		r.Get("/analytics/{short_key}", h.Analytics)
	})

	r.Get("/{short_key}", h.Redirect)

	log.Printf("url shortener listening on %s (base=%s)", addr, baseURL)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
