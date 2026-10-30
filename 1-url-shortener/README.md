# URL Shortener (Boilerplate + Judge)

A minimal, runnable Go URL shortener skeleton plus an automated correctness &
performance judge.

- `app/`   – the application (Chi router, stubbed handlers returning `501`).
- `judge/` – a standalone CLI that runs functional tests and NFR benchmarks
  against a running app instance.
- [`PROGRESS.md`](./PROGRESS.md) – step-by-step implementation checklist.

The app uses `go-chi/chi` for routing and standard SQL drivers
(`lib/pq`, `go-sql-driver/mysql`, `modernc.org/sqlite`) behind a single
`database/sql` store. The judge relies purely on the standard library.

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
│   ├── internal/storage/
│   │   ├── storage.go              # Storage interface + Link + errors
│   │   ├── sql.go                  # SQLStore (database/sql) + CRUD stubs
│   │   ├── postgres.go             # OpenPostgres (lib/pq)
│   │   ├── mysql.go                # OpenMySQL (go-sql-driver)
│   │   ├── sqlite.go               # OpenSQLite (modernc, no CGO)
│   │   └── schema.sql              # canonical table DDL
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
base. To change the port for **both** the app and the judge, set `PORT` once
when using the Makefile:

```sh
make app PORT=9090     # terminal 1
make judge PORT=9090   # terminal 2
```

`make app` passes `PORT` to both `ADDR` (`:9090`) and `BASE_URL`
(`http://localhost:9090`), and `make judge` targets the same `BASE_URL`, so the
port only needs to be changed in one place.

| Env var     | Default                      | Description                          |
| ----------- | ---------------------------- | ------------------------------------ |
| `PORT`      | `8080`                       | Port for both `ADDR` and `BASE_URL`  |
| `ADDR`      | `:<PORT>`                    | Override the listen address          |
| `BASE_URL`  | `http://localhost:<PORT>`    | Override the `short_url` prefix      |
| `DB_DRIVER` | `sqlite`                     | `sqlite` \| `postgres` \| `mysql`    |
| `DB_DSN`    | `file::memory:?cache=shared` | Driver-specific connection string    |

Example with Postgres:

```sh
DB_DRIVER=postgres \
DB_DSN="postgres://user:pass@localhost:5432/urlshortener?sslmode=disable" \
make app
```

Create the table with `app/internal/storage/schema.sql` (adjust types per
dialect as needed).

> The handlers, `SQLStore` CRUD methods, and cache methods are intentionally
> **stubbed** (`501` / not-implemented). Implement `internal/service`,
> `internal/storage/sql.go`, and `internal/cache` to make the judge pass.

## Run the judge

With the app running in another terminal:

```sh
make judge
# or
cd judge && go run . -base http://localhost:8080
```

The judge exits with code `1` when the final verdict is `FAILED`.

Every run uses a fresh random run-id for its fixtures (custom alias and long
URL), so the suite is **re-runnable against a persistent datastore** without
cleanup: no run ever collides with rows left behind by a previous run. The
trade-off is that generated rows are not deleted and will accumulate over time
(clean them up out-of-band, e.g. by resetting your dev database).

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

**Re-runnable:** every functional run uses unique fixtures (a random run-id in
the alias and long URL), so repeated runs never collide with previously
persisted test data.
