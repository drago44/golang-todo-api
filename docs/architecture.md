# Architecture Overview

This document describes the architecture and design patterns used in the Golang Todo API.

## Project Structure

```
golang-todo-api/
├── cmd/server/           # Application entry points
│   └── main.go          # Main server entry point
├── internal/            # Private application code
│   ├── app/            # Application setup and configuration
│   │   ├── config.go   # Configuration management
│   │   ├── db.go       # Database setup and connection
│   │   ├── middleware.go # HTTP middleware
│   │   └── server.go   # HTTP server setup
│   ├── router/         # HTTP routing
│   │   ├── router.go   # Route definitions
│   │   └── router_test.go
│   └── todos/          # Todo domain logic
│       ├── dto.go      # Data transfer objects
│       ├── entity.go   # Domain entities
│       ├── handlers.go # HTTP handlers
│       ├── module.go   # Dependency injection module
│       ├── repository.go # Data access layer
│       ├── service.go  # Business logic layer
│       └── *_test.go   # Unit tests
├── docs/               # Documentation
├── scripts/            # Utility scripts
└── data/              # SQLite database files
```

## Architectural Patterns

### Clean Architecture

The project follows **Clean Architecture** principles with clear separation of concerns:

1. **Entities** (`entity.go`) - Core business objects
2. **Use Cases** (`service.go`) - Business logic and rules
3. **Interface Adapters** (`handlers.go`, `repository.go`) - Data conversion and external interfaces
4. **Frameworks & Drivers** (`main.go`, `server.go`) - External frameworks and tools

### Dependency Injection

The project uses constructor-based dependency injection:

```go
// Service depends on Repository interface, not concrete implementation
type todoService struct {
    todoRepo TodoRepository
}

func NewTodoService(todoRepo TodoRepository) TodoService {
    return &todoService{todoRepo: todoRepo}
}
```

### Repository Pattern

Data access is abstracted through repository interfaces:

```go
type TodoRepository interface {
    Create(todo *Todo) (*Todo, error)
    GetAll() ([]Todo, error)
    GetByID(id uint) (*Todo, error)
    Update(todo *Todo) (*Todo, error)
    Delete(id uint) error
}
```

## Layers Description

### 1. Presentation Layer (Handlers)

**Location**: `internal/todos/handlers.go`

- Handles HTTP requests and responses
- Validates input data
- Converts between HTTP and domain models
- Uses Gin framework for HTTP handling

**Responsibilities**:
- Request parsing and validation
- Response formatting
- HTTP status code management
- Error handling and response formatting

### 2. Business Logic Layer (Services)

**Location**: `internal/todos/service.go`

- Contains all business rules and logic
- Validates business constraints
- Coordinates between different domain objects
- Independent of external frameworks

**Responsibilities**:
- Business rule enforcement
- Data validation and transformation
- Domain logic execution
- Error handling for business cases

### 3. Data Access Layer (Repository)

**Location**: `internal/todos/repository.go`

- Handles database operations
- Abstracts database-specific logic
- Implements repository interface
- Uses GORM for ORM operations

**Responsibilities**:
- Database queries and operations
- Data persistence
- Database transaction management
- Data mapping between database and domain models

### 4. Domain Layer (Entities)

**Location**: `internal/todos/entity.go`

- Defines core business objects
- Contains domain-specific validation
- Independent of external dependencies

## Technology Stack

### Core Framework
- **Gin** - HTTP web framework
- **GORM** - ORM for database operations
- **SQLite** - Database (with CGO for performance)

### Development Tools
- **Swagger** - API documentation generation
- **Testify** - Testing framework
- **golangci-lint** - Code linting
- **lefthook** - Git hooks management

### Build & Deployment
- **Docker** - Containerization with multi-stage builds
- **Make** - Build automation
- **Go modules** - Dependency management

## Design Patterns Used

### 1. Interface Segregation
Each layer defines minimal interfaces needed:
```go
type TodoService interface {
    CreateTodo(req *CreateTodoRequest) (*Todo, error)
    GetAllTodos() ([]Todo, error)
    // ... only methods this layer needs
}
```

### 2. Dependency Inversion
High-level modules don't depend on low-level modules:
```go
// Service depends on interface, not concrete repository
type todoService struct {
    todoRepo TodoRepository // interface, not concrete type
}
```

### 3. Single Responsibility
Each component has a single, well-defined responsibility:
- Handlers: HTTP concerns only
- Services: Business logic only
- Repositories: Data access only

### 4. Factory Pattern
Constructors create and wire dependencies:
```go
func NewTodoHandler(todoService TodoService) *TodoHandler
func NewTodoService(todoRepo TodoRepository) TodoService
func NewTodoRepository(db *gorm.DB) TodoRepository
```

## Data Flow

### Request Flow
1. **HTTP Request** → Gin Router
2. **Router** → Handler (Presentation Layer)
3. **Handler** → Service (Business Layer)
4. **Service** → Repository (Data Layer)
5. **Repository** → Database

### Response Flow
1. **Database** → Repository (Data Models)
2. **Repository** → Service (Domain Models)
3. **Service** → Handler (DTOs)
4. **Handler** → HTTP Response (JSON)

## Error Handling Strategy

### Domain Errors
- Defined in service layer as package-level variables
- Wrapped and handled appropriately at each layer
- Converted to appropriate HTTP status codes in handlers

```go
var (
    ErrTitleRequired = errors.New("title is required")
    ErrTitleExists   = errors.New("todo with this title already exists")
    ErrNotFound      = errors.New("todo not found")
)
```

### HTTP Error Responses
- Consistent error response format
- Appropriate HTTP status codes
- User-friendly error messages

## Configuration Management

### Environment Variables
- All configuration through environment variables
- Default values for development
- Type-safe configuration parsing

### Configuration Areas
- Server settings (host, port, mode)
- Database configuration
- Feature flags (Swagger, logging, rate limiting)
- CORS settings
- Security settings (trusted proxies)

## Testing Strategy

### Unit Tests
- Each layer tested independently
- Mock dependencies using interfaces
- High test coverage for business logic

### Integration Tests
- Router-level testing with test database
- End-to-end API testing
- Database integration testing

### Benchmarking
- Performance benchmarks for critical paths
- Memory allocation tracking
- Continuous performance monitoring

## Security Considerations

### Input Validation
- Request validation at handler level
- Business rule validation at service level
- Database constraint validation

### SQL Injection Prevention
- GORM ORM provides protection
- Prepared statements used throughout
- No raw SQL queries

### CORS Configuration
- Configurable allowed origins
- Credential handling support
- Production-ready defaults

## Performance Optimizations

### Database Optimizations
- SQLite PRAGMA settings for performance
- Prepared statement reuse
- Efficient indexing strategy
- Connection pooling

### Memory Optimizations
- Object pooling for frequently allocated objects
- Efficient JSON serialization
- Minimal memory allocations in hot paths

### Application Optimizations
- Gin release mode for production
- Static response caching
- Optimized SQL queries

## Monitoring and Observability

### Logging
- Structured logging with configurable levels
- Request/response logging
- Error logging with context

### Health Checks
- Database connectivity checks
- Application readiness indicators

### Metrics
- Performance benchmarking
- Memory usage tracking
- Request/response metrics

## Scalability Considerations

### Horizontal Scaling
- Stateless application design
- Database connection pooling
- Load balancer ready

### Vertical Scaling
- Efficient resource utilization
- Memory-optimized operations
- CPU-efficient algorithms

### Database Scaling
- Read replica support potential
- Connection pooling
- Query optimization