-- Canonical schema for the URL shortener. Apply it manually or from a
-- migration once you implement SQLStore's CRUD methods.

CREATE TABLE IF NOT EXISTS links (
    short_key    VARCHAR(32)  PRIMARY KEY,
    original_url TEXT         NOT NULL,
    click_count  BIGINT       NOT NULL DEFAULT 0,
    created_at   TIMESTAMP    NOT NULL,
    expires_at   TIMESTAMP    NULL
);

CREATE INDEX IF NOT EXISTS idx_links_expires_at ON links (expires_at);
