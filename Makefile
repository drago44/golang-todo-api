SHELL := /bin/bash

APP_NAME ?= golang-todo-api
PKG      ?= ./...
BIN_DIR  ?= bin
MAIN     ?= cmd/server/main.go

.PHONY: help test test-full-log test-short-log cover run build tidy deps fmt vet lint clean demo bench

help:
	@echo "Available targets:\n" \
	&& echo "  make test            - go test ./... -v -count=1" \
	&& echo "  make test-full-log   - richgo test ./... -v -count=1" \
	&& echo "  make test-short-log  - gotestsum --format short-verbose -- -count=1 ./..." \
	&& echo "  make cover           - coverage profile + HTML report" \
	&& echo "  make run             - run server from $(MAIN)" \
	&& echo "  make build           - build binary to $(BIN_DIR)/$(APP_NAME)" \
	&& echo "  make tidy            - go mod tidy" \
	&& echo "  make deps            - go mod download" \
	&& echo "  make fmt             - go fmt ./..." \
	&& echo "  make vet             - go vet ./..." \
	&& echo "  make lint            - golangci-lint run (if installed)" \
	&& echo "  make clean           - clean build artifacts and coverage files" \
	&& echo "  make demo            - run demo" \
	&& echo "  make bench           - run benchmarks (allocs, ns/op, B/op)"

# Tests
test:
	go test $(PKG) -v -count=1

test-full-log:
	richgo test $(PKG) -v -count=1

test-short-log:
	gotestsum --format short-verbose -- -count=1 $(PKG)

# Coverage
cover:
	go test $(PKG) -coverprofile=coverage.out -covermode=atomic
	@echo
	@echo "Coverage summary:"
	go tool cover -func=coverage.out | tail -n 1
	@echo
	@echo "Generating HTML report at coverage.html..."
	go tool cover -html=coverage.out -o coverage.html

# App lifecycle
run:
	go run $(MAIN)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN)

# Maintenance
tidy:
	go mod tidy

deps:
	go mod download

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install: brew install golangci-lint"; \
	fi

clean:
	rm -rf $(BIN_DIR) coverage.out coverage.html

demo:
	./scripts/demo.sh

bench:
	go test -bench=. -benchmem -run=^$$ ./...
