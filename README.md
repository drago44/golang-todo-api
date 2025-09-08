# Golang Todo API

A high-performance REST API for managing Todo items, built with Go, Gin, and SQLite. This project demonstrates clean architecture principles, comprehensive testing, and modern Go development practices.

## 🚀 Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd golang-todo-api

# Copy environment configuration
cp .env.example .env

# Install dependencies
make deps

# Run the application
make run
```

The API will be available at `http://localhost:8080`.

### Test the API

```bash
# Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Go", "description": "Build awesome APIs"}'

# Get all todos
curl http://localhost:8080/api/v1/todos
```

## ✨ Features

- **RESTful API** - Complete CRUD operations for todos
- **Clean Architecture** - Domain-driven design with clear separation of concerns
- **High Performance** - Optimized for speed and memory efficiency
- **Comprehensive Testing** - Unit, integration, and benchmark tests
- **Auto-documentation** - Swagger/OpenAPI integration
- **Docker Ready** - Multi-stage builds with minimal runtime footprint
- **Development Tools** - Git hooks, linting, formatting, and more
- **Production Ready** - Logging, monitoring, graceful shutdown

## 📚 Documentation

### Getting Started
- 📖 **[Development Setup](./docs/development-setup.md)** - Complete development environment setup
- 🔧 **[Configuration Guide](./docs/configuration.md)** - Environment variables and settings
- 🐳 **[Docker Guide](./docs/docker-guide.md)** - Building and running with Docker

### Architecture & Design
- 🏗️ **[Architecture Overview](./docs/architecture.md)** - System design and patterns
- 🗄️ **[Database Schema](./docs/database-schema.md)** - Database design and relationships
- 🔌 **[API Reference](./docs/api-reference.md)** - Complete API documentation

### Development
- 🧪 **[Testing Guide](./docs/testing-guide.md)** - Testing strategies and best practices
- 🛠️ **[Scripts and Tools](./docs/scripts-and-tools.md)** - Development tools and automation
- 🤝 **[Contributing Guide](./docs/contributing.md)** - How to contribute to the project

### Additional Documentation
- 📊 **[Optimization Benchmarks](./docs/optimization_benchmark.md)** - Performance improvements
- 🎣 **[Git Hooks Setup](./docs/lefthook.md)** - Automated code quality checks

## 🛠️ Technology Stack

### Core Framework
- **[Go 1.23+](https://golang.org/)** - Modern, fast, and reliable
- **[Gin](https://gin-gonic.com/)** - High-performance HTTP web framework
- **[GORM](https://gorm.io/)** - Developer-friendly ORM
- **[SQLite](https://sqlite.org/)** - Embedded database with CGO for performance

### Development Tools
- **[golangci-lint](https://golangci-lint.run/)** - Fast Go linters runner
- **[lefthook](https://github.com/evilmartians/lefthook)** - Git hooks manager
- **[gotestsum](https://github.com/gotestyourself/gotestsum)** - Enhanced test runner
- **[Swagger](https://swagger.io/)** - API documentation generation

## 🏁 Available Commands

```bash
# Development
make run              # Start the development server
make build            # Build production binary
make deps             # Download dependencies

# Code Quality  
make fmt              # Format code
make lint             # Run linters
make vet              # Run go vet

# Testing
make test             # Run all tests
make test-short       # Run short tests (faster)
make cover            # Generate coverage report
make bench            # Run benchmarks

# Documentation
make swagger-go       # Generate Swagger docs
make demo             # Run API demonstration

# Cleanup
make clean            # Remove build artifacts
```

## 📊 Performance

This API is optimized for high performance:

- **~1.47x faster** than the initial implementation
- **~1.61x less memory** usage  
- **~1.51x fewer allocations**

See [benchmark results](./docs/optimization_benchmark.md) for detailed analysis.

## 🚢 Deployment

### Using Docker

```bash
# Build and run with Docker Compose
docker compose up -d

# Or build manually
docker build -t golang-todo-api .
docker run -p 8080:8080 -v $(pwd)/data:/data golang-todo-api
```

### Production Deployment

```bash
# Build optimized binary
make build

# Configure production environment
export GIN_MODE=release
export ENABLE_SWAGGER=false
export DATABASE_URL=/var/lib/todo-api/app.db

# Run
./bin/golang-todo-api
```

## 🔧 Configuration

The application is configured via environment variables:

```bash
# Server Configuration
PORT=8080                    # Server port
HOST=localhost              # Server host
GIN_MODE=release            # Gin mode (debug/release)

# Features
ENABLE_SWAGGER=false        # Enable Swagger UI
ENABLE_LOGGER=true          # Enable request logging
ENABLE_RATE_LIMIT=false     # Enable rate limiting

# Database
DATABASE_URL=data/app.db    # SQLite database file

# CORS
ALLOWED_ORIGINS=http://localhost:3000
ALLOW_CREDENTIALS=true
```

See the [Configuration Guide](./docs/configuration.md) for complete options.

## 🧪 Testing

The project maintains high test coverage across all layers:

```bash
# Run tests with coverage
make cover

# Run benchmarks
make bench

# Run specific tests
go test ./internal/todos -v

# Race condition testing
go test -race ./...
```

## 🏗️ Project Structure

```
.
├── cmd/server/          # Application entry point
├── internal/           # Private application code
│   ├── app/           # Application setup and config
│   ├── router/        # HTTP routing
│   └── todos/         # Todo domain logic
├── docs/              # Documentation
├── scripts/           # Utility scripts  
├── data/              # SQLite database files
└── Makefile           # Build automation
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](./docs/contributing.md) for details.

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch: `git checkout -b BE-123-feature-name`
3. Make your changes with tests
4. Run quality checks: `make fmt lint test`
5. Commit your changes: `git commit -m "feat: add new feature"`
6. Push and create a Pull Request

### Development Workflow

```bash
# Install git hooks for code quality
lefthook install

# All commits are automatically:
# - Formatted with gofumpt
# - Linted with golangci-lint  
# - Tested with go test
# - Validated for conventional commits
```

## 📝 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/api/v1/todos` | List all todos |
| POST   | `/api/v1/todos` | Create new todo |
| GET    | `/api/v1/todos/{id}` | Get todo by ID |
| PUT    | `/api/v1/todos/{id}` | Update todo |
| DELETE | `/api/v1/todos/{id}` | Delete todo |

### Example Usage

```bash
# Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go",
    "description": "Master Go programming language"
  }'

# Update a todo
curl -X PUT http://localhost:8080/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn Go - Advanced",
    "completed": true
  }'

# Get all todos
curl http://localhost:8080/api/v1/todos
```

## 🐛 Troubleshooting

### Common Issues

**Port already in use:**
```bash
PORT=8081 make run
```

**Database permission errors:**
```bash
mkdir -p data
chmod 755 data
```

**Git hooks not working:**
```bash
lefthook install
```

For more troubleshooting help, see our [Development Setup Guide](./docs/development-setup.md).

## 🔍 Monitoring

### Health Check

```bash
# Basic health check
curl http://localhost:8080/api/v1/todos

# Should return: [] (empty array)
```

### Logs

The application provides structured logging:

```bash
# Enable detailed logging
ENABLE_LOGGER=true GIN_MODE=debug make run
```

## 📈 Performance Monitoring

```bash
# Run benchmarks
make bench

# Generate CPU profile
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Generate memory profile
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

## 🏆 Project Highlights

- **Clean Architecture** - Follows domain-driven design principles
- **High Test Coverage** - Comprehensive unit and integration tests
- **Performance Optimized** - Significant performance improvements documented
- **Production Ready** - Docker, logging, graceful shutdown, health checks
- **Developer Experience** - Automated formatting, linting, git hooks
- **Documentation** - Comprehensive docs for all aspects of the project

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/) - For the excellent HTTP framework
- [GORM](https://gorm.io/) - For the intuitive ORM
- [Go Team](https://golang.org/) - For the amazing programming language
- [All Contributors](./CONTRIBUTORS.md) - For making this project better

---

**Ready to build something awesome?** 🚀

Start with our [Development Setup Guide](./docs/development-setup.md) and join the community!