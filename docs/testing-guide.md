# Testing Guide

This guide covers testing strategies, tools, and best practices used in the Todo API project.

## Testing Overview

The project uses comprehensive testing approach with multiple types of tests:

- **Unit Tests** - Test individual functions and methods
- **Integration Tests** - Test component interactions
- **HTTP Tests** - Test API endpoints end-to-end
- **Benchmark Tests** - Test performance characteristics

## Testing Tools and Frameworks

### Core Testing Framework
- **Go testing** - Built-in Go testing package
- **testify** - Assertion library and test suites
- **testify/mock** - Mocking framework

### Enhanced Test Runners
- **gotestsum** - Enhanced test output and JUnit XML
- **richgo** - Colorized test output
- **go test** - Built-in test runner (fallback)

### Test Coverage
- **go test -cover** - Built-in coverage analysis
- **HTML coverage reports** - Visual coverage analysis

## Test Structure

### Test File Organization

```
internal/todos/
├── handlers.go
├── handlers_test.go      # HTTP handler tests
├── service.go
├── service_test.go       # Business logic tests
├── repository.go
├── repository_test.go    # Data access tests
├── bench_test.go         # Benchmark tests
└── module.go
```

### Test Naming Conventions

```go
// Function: CreateTodo
// Test: TestCreateTodo
// Benchmark: BenchmarkCreateTodo

func TestCreateTodo(t *testing.T) { }
func TestCreateTodo_ValidationError(t *testing.T) { }
func TestCreateTodo_DuplicateTitle(t *testing.T) { }
func BenchmarkCreateTodo(b *testing.B) { }
```

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run short tests (excluding integration tests)
make test-short

# Run tests with coverage
make cover

# Run specific package tests
go test ./internal/todos

# Run specific test function
go test ./internal/todos -run TestCreateTodo

# Verbose output
go test -v ./internal/todos
```

### Advanced Test Commands

```bash
# Run tests multiple times to catch flaky tests
go test -count=10 ./internal/todos

# Run tests in parallel
go test -parallel=4 ./...

# Race condition detection
go test -race ./...

# CPU profiling during tests
go test -cpuprofile=cpu.prof ./...

# Memory profiling during tests
go test -memprofile=mem.prof ./...
```

### Benchmark Commands

```bash
# Run all benchmarks
make bench

# Run specific benchmark
go test -bench=BenchmarkCreateTodo ./internal/todos

# Benchmark with memory statistics
go test -bench=. -benchmem ./...

# Run benchmark multiple times for accuracy
go test -bench=. -count=5 ./...
```

## Unit Testing

### Service Layer Testing

```go
func TestTodoService_CreateTodo(t *testing.T) {
    // Arrange
    mockRepo := &mockTodoRepository{}
    service := NewTodoService(mockRepo)
    
    req := &CreateTodoRequest{
        Title:       "Test Todo",
        Description: "Test Description",
    }
    
    expectedTodo := &Todo{
        ID:          1,
        Title:       req.Title,
        Description: req.Description,
        Completed:   false,
    }
    
    mockRepo.On("Create", mock.AnythingOfType("*todos.Todo")).Return(expectedTodo, nil)
    
    // Act
    result, err := service.CreateTodo(req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedTodo.Title, result.Title)
    assert.Equal(t, expectedTodo.Description, result.Description)
    assert.False(t, result.Completed)
    mockRepo.AssertExpectations(t)
}
```

### Table-Driven Tests

```go
func TestValidateTodoTitle(t *testing.T) {
    tests := []struct {
        name    string
        title   string
        wantErr bool
        errMsg  string
    }{
        {
            name:    "valid title",
            title:   "Valid Todo Title",
            wantErr: false,
        },
        {
            name:    "empty title",
            title:   "",
            wantErr: true,
            errMsg:  "title is required",
        },
        {
            name:    "whitespace only",
            title:   "   ",
            wantErr: true,
            errMsg:  "title is required",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateTodoTitle(tt.title)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Mocking with testify/mock

```go
// Mock repository implementation
type mockTodoRepository struct {
    mock.Mock
}

func (m *mockTodoRepository) Create(todo *Todo) (*Todo, error) {
    args := m.Called(todo)
    return args.Get(0).(*Todo), args.Error(1)
}

func (m *mockTodoRepository) GetAll() ([]Todo, error) {
    args := m.Called()
    return args.Get(0).([]Todo), args.Error(1)
}

// Usage in tests
mockRepo := &mockTodoRepository{}
mockRepo.On("Create", mock.AnythingOfType("*todos.Todo")).Return(&Todo{ID: 1}, nil)
```

## Integration Testing

### Database Integration Tests

```go
func TestTodoRepository_Integration(t *testing.T) {
    // Setup in-memory database
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewTodoRepository(db)
    
    t.Run("CreateAndRetrieve", func(t *testing.T) {
        // Create todo
        todo := &Todo{
            Title:       "Integration Test Todo",
            Description: "Test Description",
        }
        
        result, err := repo.Create(todo)
        assert.NoError(t, err)
        assert.NotZero(t, result.ID)
        
        // Retrieve todo
        retrieved, err := repo.GetByID(result.ID)
        assert.NoError(t, err)
        assert.Equal(t, todo.Title, retrieved.Title)
        assert.Equal(t, todo.Description, retrieved.Description)
    })
}

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    
    err = db.AutoMigrate(&Todo{})
    require.NoError(t, err)
    
    return db
}
```

## HTTP Testing

### Handler Testing

```go
func TestTodoHandler_CreateTodo(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    mockService := &mockTodoService{}
    handler := NewTodoHandler(mockService)
    
    router := gin.New()
    handler.RegisterTodoRoutes(router.Group("/api/v1"))
    
    t.Run("Success", func(t *testing.T) {
        // Arrange
        req := CreateTodoRequest{
            Title:       "Test Todo",
            Description: "Test Description",
        }
        
        expectedTodo := &Todo{
            ID:          1,
            Title:       req.Title,
            Description: req.Description,
            Completed:   false,
        }
        
        mockService.On("CreateTodo", &req).Return(expectedTodo, nil)
        
        // Prepare request
        reqBody, _ := json.Marshal(req)
        w := httptest.NewRecorder()
        httpReq, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(reqBody))
        httpReq.Header.Set("Content-Type", "application/json")
        
        // Act
        router.ServeHTTP(w, httpReq)
        
        // Assert
        assert.Equal(t, http.StatusCreated, w.Code)
        
        var response Todo
        err := json.Unmarshal(w.Body.Bytes(), &response)
        assert.NoError(t, err)
        assert.Equal(t, expectedTodo.Title, response.Title)
        
        mockService.AssertExpectations(t)
    })
}
```

### Router Integration Tests

```go
func TestRouter_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Setup dependencies
    todoRepo := NewTodoRepository(db)
    todoService := NewTodoService(todoRepo)
    
    // Setup router
    router := SetupRouter(todoService)
    
    t.Run("CreateAndGetTodo", func(t *testing.T) {
        // Create todo
        createReq := CreateTodoRequest{
            Title:       "Integration Test",
            Description: "End-to-end test",
        }
        
        reqBody, _ := json.Marshal(createReq)
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(reqBody))
        req.Header.Set("Content-Type", "application/json")
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusCreated, w.Code)
        
        var createdTodo Todo
        json.Unmarshal(w.Body.Bytes(), &createdTodo)
        
        // Get todo
        w = httptest.NewRecorder()
        req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/todos/%d", createdTodo.ID), nil)
        
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
        
        var retrievedTodo Todo
        json.Unmarshal(w.Body.Bytes(), &retrievedTodo)
        assert.Equal(t, createdTodo.ID, retrievedTodo.ID)
        assert.Equal(t, createdTodo.Title, retrievedTodo.Title)
    })
}
```

## Benchmark Testing

### Performance Benchmarks

```go
func BenchmarkCreateTodo(b *testing.B) {
    // Setup
    db := setupBenchDB(b)
    defer db.Close()
    
    repo := NewTodoRepository(db)
    service := NewTodoService(repo)
    handler := NewTodoHandler(service)
    
    router := gin.New()
    handler.RegisterTodoRoutes(router.Group("/api/v1"))
    
    req := CreateTodoRequest{
        Title:       "Benchmark Todo",
        Description: "Performance test",
    }
    
    reqBody, _ := json.Marshal(req)
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        httpReq, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(reqBody))
        httpReq.Header.Set("Content-Type", "application/json")
        
        router.ServeHTTP(w, httpReq)
        
        if w.Code != http.StatusCreated {
            b.Fatalf("Expected status %d, got %d", http.StatusCreated, w.Code)
        }
        
        // Cleanup for next iteration
        req.Title = fmt.Sprintf("Benchmark Todo %d", i)
        reqBody, _ = json.Marshal(req)
    }
}
```

### Memory Benchmark

```go
func BenchmarkTodoAllocation(b *testing.B) {
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        todo := &Todo{
            ID:          uint(i + 1),
            Title:       "Benchmark Todo",
            Description: "Memory allocation test",
            Completed:   false,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        }
        
        // Prevent compiler optimization
        _ = todo
    }
}
```

## Test Coverage

### Generating Coverage Reports

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage summary
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open coverage report
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

### Coverage Analysis

```bash
# Coverage by function
go tool cover -func=coverage.out | tail

# Find uncovered code
go tool cover -func=coverage.out | grep -E '0.0%'

# Package-level coverage
go test -cover ./...
```

### Coverage Targets

- **Overall Coverage**: > 80%
- **Handler Coverage**: > 90%
- **Service Coverage**: > 95%
- **Repository Coverage**: > 85%

## Test Configuration

### Test Environment Variables

```go
func setupTestConfig() {
    os.Setenv("GIN_MODE", "test")
    os.Setenv("DATABASE_URL", ":memory:")
    os.Setenv("ENABLE_LOGGER", "false")
}
```

### Test Database Setup

```go
func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    require.NoError(t, err)
    
    err = db.AutoMigrate(&Todo{})
    require.NoError(t, err)
    
    return db
}
```

## Continuous Integration

### GitHub Actions Testing

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22, 1.23]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: ${{ matrix.go-version }}
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: make test
    
    - name: Generate coverage
      run: make cover
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
```

### Pre-commit Testing

The project uses lefthook for pre-commit testing:

```yaml
# lefthook.yml
pre-commit:
  commands:
    tests:
      run: make test-short
      stage_fixed: true
```

## Testing Best Practices

### Writing Good Tests

1. **AAA Pattern** - Arrange, Act, Assert
2. **Descriptive Names** - Test names should describe what they test
3. **Independent Tests** - Tests should not depend on each other
4. **Fast Tests** - Unit tests should be fast
5. **Deterministic** - Tests should produce consistent results

### Test Organization

```go
func TestServiceMethod(t *testing.T) {
    // Setup common to all sub-tests
    service := setupService(t)
    
    t.Run("SuccessCase", func(t *testing.T) {
        // Test success scenario
    })
    
    t.Run("ValidationError", func(t *testing.T) {
        // Test validation error
    })
    
    t.Run("DatabaseError", func(t *testing.T) {
        // Test database error handling
    })
}
```

### Mock Best Practices

1. **Mock Interfaces** - Mock at interface boundaries
2. **Verify Calls** - Use AssertExpectations
3. **Reset Mocks** - Clean up between tests
4. **Minimal Mocks** - Only mock what you need

### Performance Testing

1. **Baseline Measurements** - Record initial performance
2. **Regression Testing** - Monitor performance changes
3. **Resource Monitoring** - Track memory and CPU usage
4. **Load Testing** - Test under realistic loads

## Troubleshooting Tests

### Common Issues

#### 1. Flaky Tests
```go
// Bad: Time-dependent test
func TestTimeout(t *testing.T) {
    time.Sleep(100 * time.Millisecond)  // Flaky
    assert.True(t, someCondition)
}

// Good: Deterministic test
func TestTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    // Use context for timeouts
}
```

#### 2. Database Connection Issues
```bash
# Error: database is locked
# Solution: Use separate test databases or :memory:
DATABASE_URL=:memory: go test
```

#### 3. Race Conditions
```bash
# Run with race detector
go test -race ./...

# Fix race conditions in code
```

### Debugging Tests

```bash
# Run single test with verbose output
go test -v -run TestSpecificFunction ./internal/todos

# Debug with delve
dlv test ./internal/todos -- -test.run TestSpecificFunction

# Print debug information
go test -v -args -debug ./internal/todos
```

### Test Performance Issues

```bash
# Profile test execution
go test -cpuprofile=test.prof ./...
go tool pprof test.prof

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# Find slow tests
go test -v ./... | grep -E "PASS.*[0-9]+\.[0-9]+s"
```

## Advanced Testing

### Fuzzing (Go 1.18+)

```go
func FuzzValidateTitle(f *testing.F) {
    // Seed corpus
    f.Add("Valid Title")
    f.Add("")
    f.Add("   ")
    
    f.Fuzz(func(t *testing.T, title string) {
        err := validateTitle(title)
        // Test that function doesn't panic
        // and returns appropriate error for invalid input
        if strings.TrimSpace(title) == "" {
            assert.Error(t, err)
        }
    })
}
```

### Contract Testing

```go
// Test that service interface is properly implemented
func TestServiceContract(t *testing.T) {
    var _ TodoService = (*todoService)(nil)  // Compile-time check
    
    // Runtime behavior testing
    service := NewTodoService(mockRepo)
    assert.Implements(t, (*TodoService)(nil), service)
}
```

### Test Utilities

```go
// Test helpers
func createTestTodo(t *testing.T, title string) *Todo {
    return &Todo{
        ID:          1,
        Title:       title,
        Description: "Test todo",
        Completed:   false,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
}

func assertTodoEqual(t *testing.T, expected, actual *Todo) {
    assert.Equal(t, expected.ID, actual.ID)
    assert.Equal(t, expected.Title, actual.Title)
    assert.Equal(t, expected.Description, actual.Description)
    assert.Equal(t, expected.Completed, actual.Completed)
}
```