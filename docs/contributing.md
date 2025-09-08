# Contributing Guide

Welcome to the Todo API project! This guide will help you contribute effectively to the project.

## Getting Started

### Prerequisites

Before contributing, ensure you have:

1. **Go 1.23+** installed
2. **Git** configured with your name and email
3. **GitHub account** for pull requests
4. **Required tools** installed (see [Development Setup](./development-setup.md))

### First Time Setup

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/golang-todo-api.git
cd golang-todo-api

# 3. Add upstream remote
git remote add upstream https://github.com/ORIGINAL_OWNER/golang-todo-api.git

# 4. Install dependencies and tools
make deps
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/evilmartians/lefthook@latest

# 5. Install git hooks
lefthook install

# 6. Verify setup
make test
```

## Development Workflow

### Branch Strategy

We use a **Git Flow** inspired branching strategy:

#### Main Branches
- **`main`** - Production-ready code
- **`dev`** - Development integration branch

#### Feature Branches
- **`BE-###-description`** - Backend features
- **`FE-###-description`** - Frontend features  
- **`FS-###-description`** - Full-stack features

#### Branch Naming Rules
```bash
# ✅ Valid branch names
git checkout -b BE-123-add-user-authentication
git checkout -b FE-456-update-todo-ui
git checkout -b FS-789-implement-search

# ❌ Invalid branch names  
git checkout -b feature-branch
git checkout -b bugfix
git checkout -b my-changes
```

### Creating a Feature Branch

```bash
# 1. Ensure you're on dev and up to date
git checkout dev
git pull upstream dev

# 2. Create feature branch
git checkout -b BE-123-add-validation

# 3. Make your changes
# ... code, test, commit ...

# 4. Push to your fork
git push origin BE-123-add-validation

# 5. Create pull request on GitHub
```

## Code Standards

### Code Style

We enforce consistent code style using automated tools:

#### Formatting
- **gofumpt** - Stricter formatting than gofmt
- **goimports** - Automatic import management
- **Applied automatically** via git hooks

#### Linting
- **golangci-lint** - Comprehensive linting
- **Configured** in `.golangci.yml`
- **Must pass** before merging

### Go Best Practices

#### Naming Conventions
```go
// ✅ Good naming
type TodoService interface { }
func NewTodoService() TodoService { }
var ErrTodoNotFound = errors.New("todo not found")

// ❌ Poor naming
type todoservice interface { }
func newTodoService() todoservice { }
var todoNotFound = errors.New("not found")
```

#### Error Handling
```go
// ✅ Proper error handling
func (s *todoService) CreateTodo(req *CreateTodoRequest) (*Todo, error) {
    if req.Title == "" {
        return nil, ErrTitleRequired
    }
    
    todo, err := s.todoRepo.Create(&Todo{
        Title:       req.Title,
        Description: req.Description,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create todo: %w", err)
    }
    
    return todo, nil
}

// ❌ Poor error handling
func (s *todoService) CreateTodo(req *CreateTodoRequest) *Todo {
    todo, _ := s.todoRepo.Create(&Todo{Title: req.Title})
    return todo  // Ignores errors
}
```

#### Interface Design
```go
// ✅ Focused interface
type TodoService interface {
    CreateTodo(req *CreateTodoRequest) (*Todo, error)
    GetTodoByID(id uint) (*Todo, error)
    UpdateTodo(id uint, req *UpdateTodoRequest) (*Todo, error)
    DeleteTodo(id uint) error
}

// ❌ Bloated interface
type TodoService interface {
    CreateTodo(req *CreateTodoRequest) (*Todo, error)
    GetTodoByID(id uint) (*Todo, error)
    // ... 20 more methods
}
```

### Testing Requirements

#### Coverage Requirements
- **Overall coverage**: >80%
- **New code**: >90%
- **Critical paths**: 100%

#### Test Types Required
```go
// 1. Unit tests for business logic
func TestTodoService_CreateTodo(t *testing.T) { }

// 2. Integration tests for data layer
func TestTodoRepository_Create(t *testing.T) { }

// 3. HTTP tests for handlers
func TestTodoHandler_CreateTodo(t *testing.T) { }

// 4. Benchmarks for performance
func BenchmarkCreateTodo(b *testing.B) { }
```

## Commit Guidelines

### Commit Message Format

We use conventional commit messages with automatic ticket prefixing:

```bash
# Your commit message
git commit -m "add user authentication endpoint"

# Automatically becomes (on branch BE-123-add-auth)
"BE-123 add user authentication endpoint"
```

#### Commit Types
- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Formatting changes
- **refactor**: Code refactoring
- **test**: Adding tests
- **chore**: Build/config changes

#### Examples
```bash
# ✅ Good commit messages
git commit -m "feat: add todo validation"
git commit -m "fix: handle database connection errors"
git commit -m "docs: update API documentation"
git commit -m "test: add integration tests for todos"

# ❌ Poor commit messages
git commit -m "fix"
git commit -m "updates"
git commit -m "WIP"
```

### Commit Best Practices

1. **Atomic commits** - One logical change per commit
2. **Clear messages** - Describe what and why, not how
3. **Small commits** - Easier to review and revert
4. **Working state** - Each commit should build and test successfully

## Pull Request Process

### Before Creating a PR

```bash
# 1. Ensure your branch is up to date
git checkout dev
git pull upstream dev
git checkout your-feature-branch
git rebase dev

# 2. Run full test suite
make test

# 3. Run linting
make lint

# 4. Format code
make fmt

# 5. Check coverage
make cover
```

### Creating a Pull Request

#### PR Title Format
```
BE-123: Add user authentication endpoint
```

#### PR Description Template
```markdown
## Summary
Brief description of changes made.

## Changes
- [ ] Added new endpoint for user authentication
- [ ] Updated database schema
- [ ] Added validation tests
- [ ] Updated documentation

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated  
- [ ] Manual testing completed
- [ ] Coverage maintained above 80%

## Breaking Changes
None / Describe any breaking changes

## Documentation
- [ ] API documentation updated
- [ ] README updated if needed
- [ ] Comments added to complex code

## Screenshots (if applicable)
[Add screenshots for UI changes]
```

### PR Review Process

#### Self-Review Checklist
- [ ] Code follows project conventions
- [ ] Tests are comprehensive and pass
- [ ] Documentation is updated
- [ ] No breaking changes (or properly documented)
- [ ] Performance impact considered
- [ ] Security implications reviewed

#### Reviewer Guidelines
- **Be constructive** - Suggest improvements, don't just point out problems
- **Be specific** - Reference line numbers and provide examples
- **Ask questions** - If something is unclear, ask for clarification
- **Approve when ready** - Don't hold up good changes for minor issues

#### Addressing Review Feedback
```bash
# Make changes based on feedback
git add .
git commit -m "address review feedback"

# Update PR
git push origin your-feature-branch
```

### Merging

- **Squash merge** - For feature branches
- **Merge commit** - For release branches
- **Delete branch** - After successful merge

## Issue Guidelines

### Creating Issues

#### Bug Reports
```markdown
**Bug Description**
Clear description of the bug.

**Steps to Reproduce**
1. Start the server
2. Send POST request to `/api/v1/todos`
3. Observe error

**Expected Behavior**
What should happen.

**Actual Behavior**  
What actually happens.

**Environment**
- OS: macOS 13.0
- Go version: 1.23
- Browser: Chrome 110

**Additional Context**
Any additional information.
```

#### Feature Requests
```markdown
**Feature Description**
Clear description of the desired feature.

**Use Case**
Why is this feature needed?

**Proposed Solution**
How should this be implemented?

**Alternatives Considered**
Other approaches you've considered.

**Additional Context**
Any additional information.
```

### Issue Labels

- **bug** - Something isn't working
- **enhancement** - New feature or request
- **documentation** - Improvements to docs
- **good first issue** - Good for newcomers
- **help wanted** - Community help needed
- **priority: high/medium/low** - Issue priority

## Code Review Guidelines

### What to Look For

#### Functionality
- Does the code do what it's supposed to do?
- Are edge cases handled?
- Is error handling appropriate?

#### Design
- Is the code well-structured?
- Does it follow project patterns?
- Is it maintainable?

#### Performance
- Are there any obvious performance issues?
- Is memory usage reasonable?
- Are database queries efficient?

#### Security
- Are there any security vulnerabilities?
- Is input validation adequate?
- Are secrets properly handled?

#### Testing
- Are there sufficient tests?
- Do tests cover edge cases?
- Are tests maintainable?

### Review Comments

#### Constructive Feedback
```
# ✅ Good feedback
Consider using a more descriptive variable name here. 
`todoCount` would be clearer than `count`.

This error handling could be improved by wrapping the error:
return fmt.Errorf("failed to create todo: %w", err)

# ❌ Poor feedback
This is wrong.
Bad variable name.
```

#### Code Suggestions
```go
// Consider this approach for better error handling:
if err != nil {
    return fmt.Errorf("failed to validate todo: %w", err)
}
```

## Documentation Standards

### Code Comments

#### When to Comment
```go
// ✅ Good comments - explain why, not what
// validateTitle checks if the title meets business requirements
// including non-empty validation and length limits
func validateTitle(title string) error { }

// calculateTax applies the current tax rate based on user location
// and product category. Rate is cached for 1 hour.
func calculateTax(amount float64, location string) float64 { }

// ❌ Unnecessary comments - code is self-explanatory
// Loop through todos
for _, todo := range todos { }

// Set completed to true  
todo.Completed = true
```

#### Documentation Comments
```go
// TodoService defines business logic for managing todos.
// All methods return domain errors that should be handled
// by the presentation layer.
type TodoService interface {
    // CreateTodo validates and creates a new todo item.
    // Returns ErrTitleRequired if title is empty.
    // Returns ErrTitleExists if title already exists.
    CreateTodo(req *CreateTodoRequest) (*Todo, error)
}
```

### API Documentation

#### Swagger Annotations
```go
// CreateTodo handles POST /todos and creates a new todo item.
// @Summary Create a new todo
// @Description Create a todo item with title and optional description
// @Tags todos
// @Accept json
// @Produce json
// @Param request body CreateTodoRequest true "Create Todo Request"
// @Success 201 {object} Todo
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 409 {object} ErrorResponse "Title already exists"
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) { }
```

## Testing Guidelines

### Test Structure

#### Arrange-Act-Assert Pattern
```go
func TestTodoService_CreateTodo(t *testing.T) {
    // Arrange
    mockRepo := &mockTodoRepository{}
    service := NewTodoService(mockRepo)
    req := &CreateTodoRequest{Title: "Test Todo"}
    
    mockRepo.On("Create", mock.AnythingOfType("*Todo")).Return(&Todo{ID: 1}, nil)
    
    // Act
    result, err := service.CreateTodo(req)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, uint(1), result.ID)
    mockRepo.AssertExpectations(t)
}
```

#### Table-Driven Tests
```go
func TestValidateTitle(t *testing.T) {
    tests := []struct {
        name    string
        title   string
        wantErr bool
    }{
        {"valid title", "Valid Todo", false},
        {"empty title", "", true},
        {"whitespace only", "   ", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateTitle(tt.title)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Performance Testing

#### Benchmarks
```go
func BenchmarkCreateTodo(b *testing.B) {
    // Setup
    service := setupBenchmarkService(b)
    req := &CreateTodoRequest{Title: "Benchmark Todo"}
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        _, err := service.CreateTodo(req)
        if err != nil {
            b.Fatal(err)
        }
        req.Title = fmt.Sprintf("Todo %d", i) // Avoid duplicates
    }
}
```

## Security Guidelines

### Input Validation
```go
// ✅ Proper validation
func (h *TodoHandler) CreateTodo(c *gin.Context) {
    var req CreateTodoRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }
    
    // Additional business validation
    if strings.TrimSpace(req.Title) == "" {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "title is required"})
        return
    }
}
```

### Error Handling
```go
// ✅ Don't expose internal details
if err != nil {
    log.Printf("Database error: %v", err) // Log details
    c.JSON(http.StatusInternalServerError, ErrorResponse{
        Error: "Internal server error", // Generic user message
    })
    return
}
```

### SQL Injection Prevention
```go
// ✅ GORM prevents SQL injection
db.Where("title = ?", userInput).Find(&todos)

// ❌ Don't use raw SQL with user input
db.Raw("SELECT * FROM todos WHERE title = " + userInput)
```

## Performance Guidelines

### Database Best Practices
```go
// ✅ Efficient queries
db.Select("id, title, completed").Where("completed = ?", false).Find(&todos)

// ❌ Inefficient queries
db.Find(&todos) // Loads all columns
```

### Memory Management
```go
// ✅ Reuse objects when possible
var todoPool = sync.Pool{
    New: func() interface{} {
        return &Todo{}
    },
}

func getTodo() *Todo {
    return todoPool.Get().(*Todo)
}

func putTodo(t *Todo) {
    t.Reset() // Clear fields
    todoPool.Put(t)
}
```

## Release Process

### Version Numbering
We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Workflow
```bash
# 1. Create release branch from dev
git checkout -b release/v1.2.0 dev

# 2. Update version and changelog
# 3. Run full test suite
make test

# 4. Create PR to main
# 5. After merge, tag release
git tag -a v1.2.0 -m "Release v1.2.0"
git push origin v1.2.0

# 6. Merge main back to dev
```

## Community Guidelines

### Code of Conduct

We are committed to providing a welcoming and inclusive environment:

1. **Be respectful** - Treat everyone with respect
2. **Be inclusive** - Welcome newcomers and different perspectives  
3. **Be collaborative** - Work together constructively
4. **Be patient** - Help others learn and grow
5. **Be professional** - Maintain professional communication

### Getting Help

- **Issues** - For bugs and feature requests
- **Discussions** - For questions and general discussion
- **Discord/Slack** - For real-time chat (if available)
- **Email** - For sensitive matters

### Recognition

We recognize contributors through:
- **Contributors file** - List of all contributors
- **Release notes** - Credit for significant contributions
- **GitHub achievements** - Automatic recognition

## FAQ

### Q: How do I run only specific tests?
```bash
# Single test function
go test -run TestCreateTodo ./internal/todos

# All tests in a package
go test ./internal/todos

# Tests matching pattern
go test -run "TestTodo.*Create" ./...
```

### Q: How do I debug failing tests?
```bash
# Verbose output
go test -v ./internal/todos

# Run with debugger
dlv test ./internal/todos -- -test.run TestCreateTodo
```

### Q: How do I add a new API endpoint?
1. Add route in `internal/router/router.go`
2. Add handler method in `internal/todos/handlers.go`
3. Add business logic in `internal/todos/service.go`
4. Add data access in `internal/todos/repository.go`
5. Add tests for all layers
6. Update Swagger documentation

### Q: How do I handle breaking changes?
1. Document the breaking change clearly
2. Provide migration guide
3. Update version number (major bump)
4. Add to changelog
5. Communicate to users

### Q: What if my PR is rejected?
- Review the feedback carefully
- Ask for clarification if needed
- Make the requested changes
- Learn from the experience
- Try again with improvements

Remember: Contributing to open source is a learning process. Don't be discouraged by feedback - use it to improve!