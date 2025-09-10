# Development Setup Guide

This guide will help you set up the development environment for the Golang Todo API project.

## Prerequisites

### Required Software

1. **Go 1.23+**
   ```bash
   # Check your Go version
   go version
   
   # If you need to install Go, visit: https://golang.org/dl/
   ```

2. **Git**
   ```bash
   git --version
   ```

3. **SQLite** (usually comes pre-installed)
   ```bash
   sqlite3 --version
   ```

### Recommended Tools

1. **golangci-lint** (for code linting)
   ```bash
   # macOS
   brew install golangci-lint
   
   # Linux/Windows
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **lefthook** (for git hooks)
   ```bash
   # Cross-platform (recommended)
   go install github.com/evilmartians/lefthook@latest
   
   # macOS (Homebrew)
   brew install lefthook
   ```

3. **gotestsum** (enhanced test runner)
   ```bash
   go install gotest.tools/gotestsum@latest
   ```

4. **richgo** (colorized test output)
   ```bash
   go install github.com/kyoh86/richgo@latest
   ```

5. **gofumpt** (stricter formatting than gofmt)
   ```bash
   go install mvdan.cc/gofumpt@latest
   ```

## Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd golang-todo-api
```

### 2. Install Dependencies

```bash
# Download Go modules
make deps

# Or manually
go mod download
```

### 3. Set Up Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your preferred settings
nano .env  # or your favorite editor
```

### 4. Set Up Git Hooks (Optional but Recommended)

```bash
# Install lefthook hooks
lefthook install
```

### 5. Run the Application

```bash
# Using make (recommended)
make run

# Or directly with go
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` by default.

### 6. Verify Setup

```bash
# Test the API
curl http://localhost:8080/api/v1/todos

# Expected response: []
```

## Development Workflow

### Daily Development Commands

```bash
# Format code
make fmt

# Run tests
make test

# Run short tests (faster)
make test-short

# Generate test coverage
make cover

# Lint code
make lint

# Vet code
make vet

# Build binary
make build

# Clean build artifacts
make clean
```

### Running with Different Configurations

```bash
# Development mode
GIN_MODE=debug make run

# Production mode
GIN_MODE=release make run

# With Swagger enabled
ENABLE_SWAGGER=true make run

# Custom port
PORT=9000 make run

# With logging enabled
ENABLE_LOGGER=true make run
```

## IDE Setup

### VS Code

**Recommended Extensions:**
- Go (official Go extension)
- REST Client (for testing API endpoints)
- GitLens (Git integration)
- Error Lens (inline error display)

**Settings** (`.vscode/settings.json`):
```json
{
    "go.formatTool": "gofumpt",
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "go.testFlags": ["-v", "-count=1"],
    "go.buildTags": "sqlite_omit_load_extension",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    }
}
```

### GoLand/IntelliJ IDEA

1. Import project as Go module
2. Enable Go modules integration
3. Set up run configurations
4. Configure code style to use gofumpt

## Environment Configuration

### Environment Variables

Copy `.env.example` to `.env` and customize:

```bash
# Server Configuration
PORT=8080                    # Server port
HOST=localhost              # Server host
PUBLIC_SCHEME=http          # Public scheme (http/https)
ENABLE_SWAGGER=false        # Enable Swagger UI
ENABLE_LOGGER=true          # Enable request logging
ENABLE_RATE_LIMIT=false     # Enable rate limiting
GIN_MODE=release           # Gin mode (debug/release)
TRUSTED_PROXIES=           # Comma-separated IPs/CIDRs

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:3000  # Allowed origins
ALLOW_CREDENTIALS=true                 # Allow credentials

# Database Configuration
DATABASE_URL=data/app.db              # SQLite database file path
```

### Configuration Validation

The application will validate configuration on startup and show helpful error messages for invalid settings.

## Database Setup

### SQLite Database

The application automatically creates the SQLite database file on first run:

```bash
# Default location
ls -la data/app.db

# Database will be created automatically when you start the server
make run
```

### Database Migrations

The application uses GORM's auto-migration feature:
- Tables are created automatically on startup
- Schema changes are applied automatically
- No manual migration scripts needed

### Database Tools

```bash
# Connect to SQLite database
sqlite3 data/app.db

# Show tables
.tables

# Show schema
.schema todos

# Exit
.exit
```

## Testing Setup

### Running Tests

```bash
# All tests
make test

# Short tests only (excluding integration tests)
make test-short

# Tests with coverage
make cover

# Specific package tests
go test ./internal/todos -v

# Specific test function
go test ./internal/todos -run TestCreateTodo -v

# Benchmarks
make bench
```

### Test Database

Tests use an in-memory SQLite database by default, so no additional setup is needed.

### Writing Tests

Follow the existing patterns in `*_test.go` files:
- Use testify for assertions
- Use table-driven tests where appropriate
- Mock external dependencies using interfaces

## Debugging

### VS Code Debugging

Create `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server/main.go",
            "env": {
                "GIN_MODE": "debug",
                "ENABLE_LOGGER": "true"
            }
        }
    ]
}
```

### Debug with Delve

```bash
# Install delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debugging session
dlv debug cmd/server/main.go
```

### Logging

Enable debug logging:
```bash
ENABLE_LOGGER=true GIN_MODE=debug make run
```

## Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Kill process using port 8080
lsof -ti:8080 | xargs kill -9

# Or use different port
PORT=8081 make run
```

#### 2. CGO Build Issues
```bash
# Ensure build tools are installed
# macOS: xcode-select --install
# Linux: apt-get install build-essential
# Windows: Install TDM-GCC or similar
```

#### 3. Module Download Issues
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
make deps
```

#### 4. Git Hooks Not Working
```bash
# Ensure lefthook is in PATH
which lefthook

# Reinstall hooks
lefthook install

# Run hooks manually
lefthook run pre-commit --all-files
```

#### 5. Database Permission Issues
```bash
# Ensure data directory exists and is writable
mkdir -p data
chmod 755 data
```

### Performance Issues

#### 1. Slow Tests
```bash
# Use short tests for development
make test-short

# Run tests in parallel
go test -parallel 4 ./...
```

#### 2. High Memory Usage
```bash
# Monitor memory usage
go test -benchmem -bench=. ./...

# Profile memory
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

## Development Best Practices

### Code Style
- Use `make fmt` before committing
- Follow Go naming conventions
- Write meaningful commit messages
- Keep functions small and focused

### Testing
- Write tests for new features
- Maintain test coverage above 80%
- Use table-driven tests for multiple scenarios
- Mock external dependencies

### Git Workflow
- Create feature branches from `dev`
- Use meaningful branch names: `BE-123-add-user-auth`
- Make atomic commits
- Write clear commit messages

### Performance
- Run benchmarks for critical paths
- Profile memory allocations
- Use `make bench` to track performance

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Documentation](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [Project Architecture](./architecture.md)
- [API Reference](./api-reference.md)
- [Contributing Guide](./contributing.md)