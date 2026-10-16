package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const nonexistentKey = "nonexistentkey123"

type functionalResult struct {
	Name   string
	Passed bool
	Detail string
	Key    string
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
	ShortKey string `json:"short_key"`
}

// judgeContext holds per-run data. All fixtures (alias, long URL) are made
// unique per run so repeated judge runs never collide with data persisted by
// earlier runs.
type judgeContext struct {
	client  *http.Client
	baseURL string
	runID   string
	alias   string
	longURL string
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		// Do not follow redirects so the judge can assert on 302/307 + Location.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func newRunID() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}

func runFunctionalTests(client *http.Client, baseURL string) []functionalResult {
	runID := newRunID()
	ctx := &judgeContext{
		client:  client,
		baseURL: baseURL,
		runID:   runID,
		alias:   "my-alias-" + runID,
		longURL: "https://example.com/very-long-path-" + runID,
	}

	tests := []func() functionalResult{
		ctx.basicFlow,
		ctx.customAlias,
		ctx.duplicateAlias,
		ctx.invalidURL,
		ctx.nonexistentKey,
	}

	results := make([]functionalResult, 0, len(tests))
	for _, fn := range tests {
		results = append(results, fn())
	}
	return results
}

func (c *judgeContext) basicFlow() functionalResult {
	res := functionalResult{Name: "Basic Shorten & Redirect"}

	status, body, err := postJSON(c.client, c.baseURL+"/api/v1/shorten",
		map[string]string{"url": c.longURL})
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if status != http.StatusCreated {
		res.Detail = fmt.Sprintf("expected 201, got %d", status)
		return res
	}

	var sr shortenResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		res.Detail = "invalid JSON response: " + err.Error()
		return res
	}
	if sr.ShortKey == "" {
		res.Detail = "empty short_key in response"
		return res
	}
	res.Key = sr.ShortKey

	loc, rstatus, err := getRedirect(c.client, c.baseURL+"/"+sr.ShortKey)
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if rstatus != http.StatusFound && rstatus != http.StatusTemporaryRedirect {
		res.Detail = fmt.Sprintf("expected 302/307, got %d", rstatus)
		return res
	}
	if loc != c.longURL {
		res.Detail = fmt.Sprintf("Location = %q, want %q", loc, c.longURL)
		return res
	}

	res.Passed = true
	return res
}

func (c *judgeContext) customAlias() functionalResult {
	res := functionalResult{Name: "Custom Alias Support"}

	status, body, err := postJSON(c.client, c.baseURL+"/api/v1/shorten",
		map[string]string{"url": "https://example.com/custom-alias-" + c.runID, "custom_alias": c.alias})
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if status != http.StatusCreated {
		res.Detail = fmt.Sprintf("expected 201, got %d", status)
		return res
	}

	var sr shortenResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		res.Detail = "invalid JSON response: " + err.Error()
		return res
	}
	if sr.ShortKey != c.alias {
		res.Detail = fmt.Sprintf("short_key = %q, want %q", sr.ShortKey, c.alias)
		return res
	}
	res.Key = sr.ShortKey
	res.Passed = true
	return res
}

func (c *judgeContext) duplicateAlias() functionalResult {
	res := functionalResult{Name: "Duplicate Alias Conflict (409)"}

	status, _, err := postJSON(c.client, c.baseURL+"/api/v1/shorten",
		map[string]string{"url": "https://example.com/duplicate-" + c.runID, "custom_alias": c.alias})
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if status != http.StatusConflict {
		res.Detail = fmt.Sprintf("expected 409, got %d", status)
		return res
	}
	res.Passed = true
	return res
}

func (c *judgeContext) invalidURL() functionalResult {
	res := functionalResult{Name: "Invalid URL Validation (400)"}

	status, _, err := postJSON(c.client, c.baseURL+"/api/v1/shorten",
		map[string]string{"url": "not-a-url"})
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if status != http.StatusBadRequest {
		res.Detail = fmt.Sprintf("expected 400, got %d", status)
		return res
	}
	res.Passed = true
	return res
}

func (c *judgeContext) nonexistentKey() functionalResult {
	res := functionalResult{Name: "Non-Existent Key Handling (404)"}

	_, status, err := getRedirect(c.client, c.baseURL+"/"+nonexistentKey)
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if status != http.StatusNotFound {
		res.Detail = fmt.Sprintf("expected 404, got %d", status)
		return res
	}
	res.Passed = true
	return res
}

func postJSON(client *http.Client, url string, payload any) (int, []byte, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, body, nil
}

func getRedirect(client *http.Client, url string) (string, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.Header.Get("Location"), resp.StatusCode, nil
}
