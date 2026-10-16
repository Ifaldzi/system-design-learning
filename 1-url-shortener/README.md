# URL Shortener (Boilerplate + Judge)

A minimal, runnable Go URL shortener skeleton plus an automated correctness &
performance judge.

- `app/`   – the application (Chi router, stubbed handlers returning `501`).
- `judge/` – a standalone CLI that runs functional tests and NFR benchmarks
  against a running app instance.
- [`PROGRESS.md`](./PROGRESS.md) – step-by-step implementation checklist.

Only one external dependency is used: `github.com/go-chi/chi/v5`.
The judge relies purely on the standard library.

## Requirements

- Go 1.22+

## Layout

```
1-url-shortener/
├── app/
│   ├── cmd/server/main.go          # startup + Chi wireup
│   ├── internal/handler/
│   │   ├── shortener.go            # HTTP handler stubs (501)
│   │   └── health.go               # GET /health
│   ├── internal/service/service.go # core logic interfaces/stubs
│   ├── internal/storage/storage.go # persistence interface + memory stub
│   └── internal/cache/cache.go     # read-through cache interface + memory stub
├── judge/
│   ├── main.go                     # CLI entrypoint / flags
│   ├── functional.go               # functional API tests
│   ├── benchmark.go                # load tests + percentiles
│   └── report.go                   # terminal dashboard
├── Makefile
└── README.md
```

## API contract

| Method | Path                          | Success                  | Errors                          |
| ------ | ----------------------------- | ------------------------ | ------------------------------- |
| POST   | `/api/v1/shorten`             | `201` `{short_url,short_key}` | `400` invalid, `409` alias taken |
| GET    | `/{short_key}`                | `302`/`307` + `Location` | `404` unknown/expired           |
| GET    | `/api/v1/analytics/{short_key}` | `200` metadata         | `404` unknown                   |
| GET    | `/health`                     | `200` `{"status":"ok"}`  | –                               |

Request body for `POST /api/v1/shorten`:

```json
{ "url": "https://example.com/very-long-path", "custom_alias": "optional-alias" }
```

## Run the app

```sh
make app
# or
cd app && go run ./cmd/server
```

The server listens on `:8080` and uses `http://localhost:8080` as the short-URL
base. Both are overridable with `ADDR` / `BASE_URL` environment variables.

> The handlers are intentionally **stubbed** and return `501 Not Implemented`.
> Implement the `internal/service` and `internal/storage` logic to make the
> judge pass.

## Run the judge

With the app running in another terminal:

```sh
make judge
# or
cd judge && go run . -base http://localhost:8080
```

The judge exits with code `1` when the final verdict is `FAILED`.

### Judge flags

| Flag                     | Default                 | Description                          |
| ------------------------ | ----------------------- | ------------------------------------ |
| `-base`                  | `http://localhost:8080` | Target base URL                      |
| `-n`                     | `10000`                 | Requests per stress test             |
| `-c`                     | `50`                    | Concurrent workers                   |
| `-penetration-rps`       | `1000`                  | Cache penetration request rate       |
| `-penetration-duration`  | `3s`                    | Cache penetration duration           |
| `-read-p95-target`       | `10ms`                  | Read p95 latency target              |
| `-write-p95-target`      | `50ms`                  | Write p95 latency target             |

## What the judge checks

**Functional:** basic shorten→redirect, custom alias, duplicate alias `409`,
invalid URL `400`, unknown key `404`.

**NFR:** read latency p50/p95/p99 under concurrent load, write latency p95,
a 10k-request stress test across 50 workers (no `500`s), and a cache-penetration
test firing 1,000 rps at non-existent keys.
