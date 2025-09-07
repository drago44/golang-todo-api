package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drago44/golang-todo-api/internal/todos"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service to use with real TodoHandler for route wiring
type mockService struct{ mock.Mock }

func (m *mockService) CreateTodo(req *todos.CreateTodoRequest) (*todos.Todo, error) {
	args := m.Called(req)
	if v := args.Get(0); v != nil {
		return v.(*todos.Todo), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockService) GetAllTodos() ([]todos.Todo, error) {
	args := m.Called()
	return args.Get(0).([]todos.Todo), args.Error(1)
}

func (m *mockService) GetTodoByID(id uint) (*todos.Todo, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*todos.Todo), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockService) UpdateTodo(id uint, req *todos.UpdateTodoRequest) (*todos.Todo, error) {
	args := m.Called(id, req)
	if v := args.Get(0); v != nil {
		return v.(*todos.Todo), args.Error(1)
	}

	return nil, args.Error(1)
}
func (m *mockService) DeleteTodo(id uint) error { return m.Called(id).Error(0) }

func TestRouter_HealthAndTodosRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	mockSvc := new(mockService)
	// For GET /api/v1/todos
	mockSvc.On("GetAllTodos").Return([]todos.Todo{}, nil).Once()
	h := todos.NewTodoHandler(mockSvc)

	r := New(engine, h, false)

	// Health
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.GetEngine().ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	t.Logf("GET /health status=%d body=%s", w.Code, w.Body.String())

	// GET /api/v1/todos
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	w2 := httptest.NewRecorder()
	r.GetEngine().ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	t.Logf("GET /api/v1/todos status=%d body=%s", w2.Code, w2.Body.String())

	// Ensure routes are registered
	var hasHealth, hasTodos bool

	for _, ri := range r.GetEngine().Routes() {
		if ri.Path == "/health" && ri.Method == http.MethodGet {
			hasHealth = true
		}

		if ri.Path == "/api/v1/todos" && ri.Method == http.MethodGet {
			hasTodos = true
		}
	}

	assert.True(t, hasHealth)
	assert.True(t, hasTodos)

	mockSvc.AssertExpectations(t)
}
