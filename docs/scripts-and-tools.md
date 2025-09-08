# Scripts and Tools Documentation

This document describes all the scripts, tools, and automation used in the project.

## Project Scripts

### Demo Script (`scripts/demo.sh`)

The demo script showcases the API functionality with automated testing and demonstration.

#### Usage

```bash
# Run full demo (tests + API demo)
./scripts/demo.sh

# Run only tests
./scripts/demo.sh test

# Run only API demonstration
./scripts/demo.sh api

# Show usage help
./scripts/demo.sh --help
```

#### What it does

1. **Environment Setup**
   - Creates temporary directory (`.tmp/`)
   - Sets up isolated SQLite database
   - Configures demo environment variables

2. **API Testing** (when `test` or `demo` mode)
   - Starts the API server in background
   - Waits for server to be ready
   - Runs comprehensive API tests
   - Demonstrates all CRUD operations
   - Shows error handling scenarios

3. **API Demonstration** (when `api` or `demo` mode)
   - Interactive API examples
   - Creates sample todos
   - Shows various API endpoints
   - Demonstrates data persistence

#### Configuration

```bash
# Environment variables for demo
HOST="${HOST:-localhost}"          # Server host
PORT="${PORT:-8080}"              # Server port
BASE_URL="http://${HOST}:${PORT}" # Base API URL
DEMO_DIR="${ROOT_DIR}/.tmp"       # Demo directory
DB_PATH="${DEMO_DIR}/app.db"      # Demo database
DEBUG_LOG="${DEBUG_LOG:-false}"   # Enable debug logging
```

#### Example Output

```bash
$ ./scripts/demo.sh

=== Todo API Demo ===
🚀 Starting API server on http://localhost:8080
✅ Server is ready!

📝 Creating todos...
✅ Created: "Buy groceries" (ID: 1)
✅ Created: "Walk the dog" (ID: 2)
✅ Created: "Finish project" (ID: 3)

📋 Getting all todos...
Found 3 todos:
- [1] Buy groceries (pending)
- [2] Walk the dog (pending) 
- [3] Finish project (pending)

✏️ Updating todo...
✅ Updated todo 2: "Walk the dog" → completed

🗑️ Deleting todo...
✅ Deleted todo 3

📊 Final state: 2 todos remaining

🧹 Cleaning up...
✅ Demo completed successfully!
```

## Makefile Commands

The project uses a comprehensive Makefile for build automation and development tasks.

### Available Commands

#### Development Commands

```bash
# Show all available commands
make help

# Run the application
make run

# Build binary
make build

# Clean build artifacts
make clean
```

#### Testing Commands

```bash
# Run all tests
make test

# Run short tests (faster)
make test-short

# Generate coverage report
make cover

# Run benchmarks
make bench
```

#### Code Quality Commands

```bash
# Format code
make fmt

# Vet code
make vet

# Lint code
make lint

# Download dependencies
make deps

# Tidy go.mod
make tidy
```

#### Swagger Commands

```bash
# Generate Swagger docs (Go only)
make swagger-go

# Generate all Swagger formats (Go, JSON, YAML)
make swagger-full
```

### Makefile Configuration

#### Variables

```makefile
APP_NAME ?= golang-todo-api    # Application name
PKG      ?= ./...             # Package pattern
BIN_DIR  ?= bin              # Binary output directory
MAIN     ?= cmd/server/main.go # Main entry point
GOPATH_BIN := $(shell go env GOPATH)/bin
SWAG_RUN   := go run github.com/swaggo/swag/cmd/swag@v1.16.6
```

#### Tool Detection

The Makefile automatically detects and uses enhanced tools when available:

```makefile
# Enhanced test runner
test-short:
	@if command -v gotestsum >/dev/null 2>&1; then \
		gotestsum --format short-verbose -- -count=1 $(PKG); \
	else \
		go test $(PKG) -v -count=1; \
	fi

# Better formatting
fmt:
	@if command -v gofumpt >/dev/null 2>&1; then \
		gofumpt -w ./; \
	else \
		go fmt ./...; \
	fi
```

## Git Hooks (Lefthook)

### Configuration (`lefthook.yml`)

The project uses lefthook for automated Git hooks that ensure code quality.

#### Pre-commit Hooks

```yaml
pre-commit:
  commands:
    branch-name:
      run: scripts/lefthook/git-validate-branch.sh
      
    fmt-staged:
      run: scripts/lefthook/fmt-staged.sh
      stage_fixed: true
      
    vet:
      run: make vet
      
    test-short:
      run: make test-short
      
    lint:
      run: make lint
```

#### Commit Message Hook

```yaml
commit-msg:
  commands:
    validate-commit-msg:
      run: scripts/lefthook/git-validate-commit-msg.sh {1}
```

#### Pre-push Hook

```yaml
pre-push:
  commands:
    test:
      run: make test-short
```

### Lefthook Scripts

#### Branch Name Validation (`scripts/lefthook/git-validate-branch.sh`)

Validates branch names against allowed patterns:

- `dev` - Development branch
- `main` - Main branch  
- `BE-####-description` - Backend feature branches
- `FE-####-description` - Frontend feature branches
- `FS-####-description` - Full-stack feature branches

```bash
# Valid branch names
git checkout -b BE-123-add-user-auth
git checkout -b FE-456-update-ui
git checkout -b FS-789-new-feature

# Invalid branch names (will be rejected)
git checkout -b feature-branch
git checkout -b bugfix
```

#### Staged File Formatting (`scripts/lefthook/fmt-staged.sh`)

Automatically formats staged Go files:

1. Identifies staged `.go` files
2. Formats using `gofumpt` (or `gofmt` as fallback)
3. Re-adds formatted files to staging area
4. Prevents empty commits if no changes remain

#### Commit Message Validation (`scripts/lefthook/git-validate-commit-msg.sh`)

Automatically prepends ticket prefixes to commit messages:

```bash
# Branch: BE-123-add-auth
# Commit: "add user authentication"
# Result: "BE-123 add user authentication"

# Skip for special commits
git commit -m "Merge branch 'feature'"    # Not modified
git commit -m "Revert previous commit"    # Not modified
```

### Installing and Managing Hooks

```bash
# Install hooks
lefthook install

# Run hooks manually
lefthook run pre-commit --all-files
lefthook run commit-msg .git/COMMIT_EDITMSG

# Uninstall hooks
lefthook uninstall

# Skip hooks temporarily
LEFTHOOK=0 git commit -m "skip hooks"
```

## Development Tools

### Code Quality Tools

#### golangci-lint

Configuration in `.golangci.yml`:

```yaml
run:
  timeout: 5m
  tests: false

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - misspell
    - gofmt
    - goimports
```

Usage:
```bash
# Install
brew install golangci-lint

# Run linting
make lint
golangci-lint run
```

#### gofumpt (Enhanced formatting)

```bash
# Install
go install mvdan.cc/gofumpt@latest

# Format code
make fmt
gofumpt -w .
```

### Testing Tools

#### gotestsum (Enhanced test runner)

```bash
# Install
go install gotest.tools/gotestsum@latest

# Run tests
make test-short
gotestsum --format short-verbose
```

#### richgo (Colorized test output)

```bash
# Install  
go install github.com/kyoh86/richgo@latest

# Run tests
make test
richgo test ./...
```

### Documentation Tools

#### Swagger/OpenAPI

```bash
# Generate documentation
make swagger-go      # Go bindings only
make swagger-full    # All formats

# The swagger command used
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init \
  -g cmd/server/main.go \
  -o docs/swagger \
  -d . \
  --outputTypes go,json,yaml \
  --parseInternal
```

## CI/CD Integration

### GitHub Actions

The tools integrate with GitHub Actions workflows:

```yaml
- name: Install tools
  run: |
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    go install gotest.tools/gotestsum@latest

- name: Lint
  run: make lint

- name: Test
  run: make test-short

- name: Build
  run: make build
```

### Docker Integration

Tools are available in Docker containers:

```dockerfile
# Install development tools
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install gotest.tools/gotestsum@latest && \
    go install mvdan.cc/gofumpt@latest
```

## IDE Integration

### VS Code

#### Recommended Extensions

```json
{
    "recommendations": [
        "golang.go",
        "ms-vscode.vscode-json",
        "timonwong.shellcheck",
        "esbenp.prettier-vscode"
    ]
}
```

#### Settings

```json
{
    "go.formatTool": "gofumpt",
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "go.testFlags": ["-v", "-count=1"],
    "editor.formatOnSave": true
}
```

### GoLand/IntelliJ IDEA

1. **File Watchers**: Set up for gofumpt formatting
2. **Go Linter**: Configure golangci-lint integration
3. **Run Configurations**: Create configs for make commands

## Troubleshooting Tools

### Common Issues

#### Tool Not Found

```bash
# Error: command not found
# Solution: Install tool or check PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
which golangci-lint
which gotestsum
```

#### Slow Linting

```bash
# Use faster linting configuration
golangci-lint run --fast

# Or disable slow linters in .golangci.yml
```

#### Git Hooks Not Running

```bash
# Check lefthook installation
lefthook version

# Reinstall hooks
lefthook uninstall
lefthook install

# Check hook permissions
ls -la .git/hooks/
```

#### Demo Script Issues

```bash
# Port already in use
PORT=8081 ./scripts/demo.sh

# Permission issues
chmod +x scripts/demo.sh

# Database issues
rm -rf .tmp/
./scripts/demo.sh
```

### Debugging Tools

#### Verbose Makefile

```bash
# See all commands being executed
make -n test    # Dry run
make V=1 test   # Verbose output
```

#### Git Hook Debugging

```bash
# Run individual hooks
.git/hooks/pre-commit
.git/hooks/commit-msg .git/COMMIT_EDITMSG

# Debug lefthook
LEFTHOOK_VERBOSE=1 git commit
```

#### Demo Script Debugging

```bash
# Enable debug output
DEBUG_LOG=true ./scripts/demo.sh

# Manual testing
export PORT=8081
export DEBUG_LOG=true
./scripts/demo.sh api
```

## Custom Tool Configuration

### Adding New Tools

#### 1. Update Makefile

```makefile
new-tool:
	@if command -v newtool >/dev/null 2>&1; then \
		newtool run; \
	else \
		echo "newtool not installed. Install: go install example.com/newtool@latest"; \
	fi
```

#### 2. Add to Lefthook

```yaml
pre-commit:
  commands:
    new-check:
      run: make new-tool
```

#### 3. Update Documentation

```markdown
### NewTool
- **Purpose**: Description of tool
- **Installation**: `go install example.com/newtool@latest`
- **Usage**: `make new-tool`
```

### Tool Version Management

```bash
# Pin tool versions in go.mod (tools.go)
//go:build tools
// +build tools

package tools

import (
    _ "github.com/golangci/golangci-lint/cmd/golangci-lint"
    _ "gotest.tools/gotestsum"
)
```

```bash
# Install pinned versions
go install github.com/golangci/golangci-lint/cmd/golangci-lint
go install gotest.tools/gotestsum
```

## Performance Monitoring

### Benchmark Integration

```bash
# Run benchmarks and save results
make bench | tee benchmark-$(date +%Y%m%d).txt

# Compare benchmarks
benchcmp old.txt new.txt

# Continuous benchmarking
git log --oneline | head -10 | while read commit; do
    git checkout $commit
    make bench > bench-$commit.txt
done
```

### Profiling Tools

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling  
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Block profiling
go test -blockprofile=block.prof -bench=.
go tool pprof block.prof
```