APP_DIR   := app
JUDGE_DIR := judge
PORT      ?= 8080
BASE_URL  ?= http://localhost:$(PORT)

.PHONY: help app judge build tidy

help:
	@echo "URL Shortener"
	@echo ""
	@echo "  make app     - run the shortener server on :$(PORT)"
	@echo "  make judge   - run the judge against BASE_URL ($(BASE_URL))"
	@echo "  make build   - build the app and judge modules"
	@echo "  make tidy    - go mod tidy for both modules"
	@echo ""
	@echo "  Override the port for both app and judge with: make <target> PORT=9090"

app:
	cd $(APP_DIR) && PORT=$(PORT) go run ./cmd/server

judge:
	cd $(JUDGE_DIR) && go run . -base $(BASE_URL)

build:
	cd $(APP_DIR) && go build ./...
	cd $(JUDGE_DIR) && go build ./...

tidy:
	cd $(APP_DIR) && go mod tidy
	cd $(JUDGE_DIR) && go mod tidy
