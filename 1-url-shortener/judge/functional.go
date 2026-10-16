package main

import (
	"bytes"
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

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		// Do not follow redirects so the judge can assert on 302/307 + Location.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func runFunctionalTests(client *http.Client, baseURL string) []functionalResult {
	tests := []func(*http.Client, string) functionalResult{
		testBasicFlow,
		testCustomAlias,
		testDuplicateAlias,
		testInvalidURL,
		testNonexistentKey,
	}

	results := make([]functionalResult, 0, len(tests))
	for _, fn := range tests {
		results = append(results, fn(client, baseURL))
	}
	return results
}

func testBasicFlow(client *http.Client, baseURL string) functionalResult {
	res := functionalResult{Name: "Basic Shorten & Redirect"}
	const longURL = "https://example.com/very-long-path"

	status, body, err := postJSON(client, baseURL+"/api/v1/shorten",
		map[string]string{"url": longURL})
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

	loc, rstatus, err := getRedirect(client, baseURL+"/"+sr.ShortKey)
	if err != nil {
		res.Detail = err.Error()
		return res
	}
	if rstatus != http.StatusFound && rstatus != http.StatusTemporaryRedirect {
		res.Detail = fmt.Sprintf("expected 302/307, got %d", rstatus)
		return res
	}
	if loc != longURL {
		res.Detail = fmt.Sprintf("Location = %q, want %q", loc, longURL)
		return res
	}

	res.Passed = true
	return res
}

func testCustomAlias(client *http.Client, baseURL string) functionalResult {
	res := functionalResult{Name: "Custom Alias Support"}
	const alias = "my-alias"

	status, body, err := postJSON(client, baseURL+"/api/v1/shorten",
		map[string]string{"url": "https://example.com/custom-alias", "custom_alias": alias})
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
	if sr.ShortKey != alias {
		res.Detail = fmt.Sprintf("short_key = %q, want %q", sr.ShortKey, alias)
		return res
	}
	res.Key = sr.ShortKey
	res.Passed = true
	return res
}

func testDuplicateAlias(client *http.Client, baseURL string) functionalResult {
	res := functionalResult{Name: "Duplicate Alias Conflict (409)"}
	const alias = "my-alias"

	status, _, err := postJSON(client, baseURL+"/api/v1/shorten",
		map[string]string{"url": "https://example.com/duplicate", "custom_alias": alias})
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

func testInvalidURL(client *http.Client, baseURL string) functionalResult {
	res := functionalResult{Name: "Invalid URL Validation (400)"}

	status, _, err := postJSON(client, baseURL+"/api/v1/shorten",
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

func testNonexistentKey(client *http.Client, baseURL string) functionalResult {
	res := functionalResult{Name: "Non-Existent Key Handling (404)"}

	_, status, err := getRedirect(client, baseURL+"/"+nonexistentKey)
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
