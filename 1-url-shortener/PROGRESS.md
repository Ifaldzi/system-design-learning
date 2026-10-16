# URL Shortener — Learning Progress Checklist

Work top to bottom. Each phase has a goal, the files you'll touch, and the judge
check that proves it's done. Tick items with `x` as you go.

Run the server: `make app`  
Run the judge:  `make judge`

---

## Phase 0 — Understand the baseline

- [ ] Read `app/cmd/server/main.go` and trace every route to its handler.
- [ ] Read `app/internal/service/service.go` and `app/internal/storage/storage.go`.
- [ ] Start the server and confirm `curl localhost:8080/health` returns `{"status":"ok"}`.
- [ ] Confirm every stubbed handler currently returns `501 Not Implemented`.
- [ ] Run the judge once and read the failing report (baseline snapshot).

**Done when:** you can explain request flow `router -> handler -> service -> storage`.

---

## Phase 1 — Implement storage (foundation)

- [ ] `storage.Memory.Save` — store a `*Link`, return `ErrConflict` on duplicate key.
- [ ] `storage.Memory.Get` — return `ErrNotFound` for missing/expired links.
- [ ] `storage.Memory.IncrementClicks` — atomically bump `ClickCount` (`sync.Mutex`).
- [ ] Add a helper to generate a random short key (e.g. base62 of 6 chars).
- [ ] Sanity-test storage with a small Go test (`go test ./...`).

**Done when:** storage round-trips a link and blocks duplicate keys.

---

## Phase 2 — `POST /api/v1/shorten` (core write path)

- [ ] Decode JSON body `{"url", "custom_alias"}` in `handler.Create`.
- [ ] Validate `url` — parse + require `http`/`https` scheme; else `400`.
- [ ] If `custom_alias` empty, generate a key; else validate alias format.
- [ ] Call `service.Shorten`, persist via storage, return `201`.
- [ ] Respond `{"short_url": "<baseURL>/<key>", "short_key": "<key>"}`.
- [ ] Map errors: invalid input -> `400`, alias taken -> `409`.

**Judge checks:** `Invalid URL Validation (400)`, `Custom Alias Support`.

---

## Phase 3 — `GET /{short_key}` (core read path)

- [ ] Read `chi.URLParam(r, "short_key")` in `handler.Redirect`.
- [ ] Call `service.Resolve`; `404` if unknown or expired.
- [ ] On success, increment clicks then `http.Redirect` (`302`/`307`) with the original URL.
- [ ] Verify the `Location` header is exactly the stored long URL.

**Judge checks:** `Basic Shorten & Redirect`, `Non-Existent Key Handling (404)`.

---

## Phase 4 — Custom alias collision handling

- [ ] Ensure `Save` returns `ErrConflict` when the key already exists.
- [ ] Ensure `service.Shorten` propagates it as `ErrAliasTaken`.
- [ ] Ensure `handler.Create` maps it to `409 Conflict`.

**Judge checks:** `Duplicate Alias Conflict (409)`.

---

## Phase 5 — `GET /api/v1/analytics/{short_key}`

- [ ] Add `handler.Analytics`: read key, call `service.Analytics`.
- [ ] Return `200` with `{"short_key","original_url","click_count","created_at"}`.
- [ ] Return `404` for unknown keys.

**Done when:** click counts match the number of successful redirects.

---

## Phase 6 — Expiration

- [ ] Set `ExpiresAt` on new links (e.g. TTL configurable).
- [ ] `Get` must treat expired links as `ErrNotFound` (-> `404`).

**Done when:** an expired key stops redirecting.

---

## Phase 7 — Non-functional requirements

- [ ] Confirm no data races: `cd app && go test -race ./...`.
- [ ] Keep read path fast — reader lock only, no blocking on the hot path.
- [ ] Cache negative lookups (or make misses cheap) for penetration immunity.
- [ ] Verify no handler ever returns `500` under load.

**Judge checks:** read p95 `<= 10ms`, write p95 `<= 50ms`, cache-penetration p95
`<= 10ms`, zero `500`s across 10k requests / 50 workers.

---

## Phase 8 — Green run

- [ ] All 5 functional tests `[PASS]`.
- [ ] `FINAL VERDICT: PASSED` from `make judge`.
- [ ] `go vet ./...` and `gofmt -l .` are clean.
- [ ] Commit the working implementation.

---

## Stretch goals (after passing)

- [ ] Persist to a real DB (SQLite/Postgres) behind the `Storage` interface.
- [ ] Add rate limiting per IP.
- [ ] Add structured logging / metrics (p50/p95/p99 export).
- [ ] Add graceful shutdown and configurable TTL/base62 length.
