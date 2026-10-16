APP_DIR   := app
JUDGE_DIR := judge
BASE_URL  ?= http://localhost:8080

.PHONY: help app judge build tidy

help:
	@echo "URL Shortener"
	@echo ""
	@echo "  make app     - run the shortener server on :8080"
	@echo "  make judge   - run the judge against BASE_URL ($(BASE_URL))"
	@echo "  make build   - build the app and judge modules"
	@echo "  make tidy    - go mod tidy for both modules"

app:
	cd $(APP_DIR) && go run ./cmd/server

judge:
	cd $(JUDGE_DIR) && go run . -base $(BASE_URL)

build:
	cd $(APP_DIR) && go build ./...
	cd $(JUDGE_DIR) && go build ./...

tidy:
	cd $(APP_DIR) && go mod tidy
	cd $(JUDGE_DIR) && go mod tidy
