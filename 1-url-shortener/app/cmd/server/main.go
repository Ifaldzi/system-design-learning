package main

import (
	"context"
	"fmt"
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
	port := env("PORT", "8080")
	baseURL := env("BASE_URL", "http://localhost:"+port)
	addr := env("ADDR", ":"+port)

	store, err := openStore(context.Background())
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer store.Close()

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

// openStore selects the SQL backend from env:
//
//	DB_DRIVER=sqlite|postgres|mysql   (default: sqlite)
//	DB_DSN=<driver specific DSN>      (default: shared in-memory sqlite)
func openStore(ctx context.Context) (*storage.SQLStore, error) {
	driver := env("DB_DRIVER", "sqlite")
	dsn := env("DB_DSN", "file::memory:?cache=shared")

	switch driver {
	case "postgres":
		return storage.OpenPostgres(ctx, dsn)
	case "mysql":
		return storage.OpenMySQL(ctx, dsn)
	case "sqlite":
		return storage.OpenSQLite(ctx, dsn)
	default:
		return nil, fmt.Errorf("unknown DB_DRIVER %q (want sqlite|postgres|mysql)", driver)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
