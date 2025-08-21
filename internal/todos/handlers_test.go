package todos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service for handler tests
type mockTodoService struct{ mock.Mock }

func (m *mockTodoService) CreateTodo(req *CreateTodoRequest) (*Todo, error) {
	args := m.Called(req)
	if v := args.Get(0); v != nil {
		return v.(*Todo), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTodoService) GetAllTodos() ([]Todo, error) {
	args := m.Called()
	return args.Get(0).([]Todo), args.Error(1)
}
func (m *mockTodoService) GetTodoByID(id uint) (*Todo, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*Todo), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTodoService) UpdateTodo(id uint, req *UpdateTodoRequest) (*Todo, error) {
	args := m.Called(id, req)
	if v := args.Get(0); v != nil {
		return v.(*Todo), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockTodoService) DeleteTodo(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupRouter(handler *TodoHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	rg := r.Group("/")
	handler.RegisterTodoRoutes(rg)
	return r
}

func TestCreateTodo_Success_Handler(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	body := CreateTodoRequest{Title: "A", Description: "d"}
	b, _ := json.Marshal(body)
	t.Logf("HTTP POST /todos: body=%s", string(b))

	mockSvc.On("CreateTodo", &body).Return(&Todo{ID: 1, Title: "A", Description: "d"}, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/todos/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP POST /todos: status=%d resp=%s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp Todo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "A", resp.Title)

	mockSvc.AssertExpectations(t)
}

func TestCreateTodo_BadRequest(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	// Missing required title
	b := []byte(`{"description":"d"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP POST /todos (bad): status=%d resp=%s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAllTodos_Success(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	expected := []Todo{{ID: 1, Title: "A"}}
	mockSvc.On("GetAllTodos").Return(expected, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/todos/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP GET /todos: status=%d items", w.Code)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []Todo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Len(t, resp, 1)
	assert.Equal(t, "A", resp[0].Title)

	mockSvc.AssertExpectations(t)
}

func TestGetTodoByID_InvalidID(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/todos/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP GET /todos/abc: status=%d resp=%s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertNotCalled(t, "GetTodoByID", mock.Anything)
}

func TestGetTodoByID_NotFound(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	mockSvc.On("GetTodoByID", uint(2)).Return(nil, assert.AnError).Once()

	req := httptest.NewRequest(http.MethodGet, "/todos/2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP GET /todos/2: status=%d resp=%s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateTodo_Success_Handler(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	completed := true
	// Title is required by binding; keep same title to avoid uniqueness
	body := UpdateTodoRequest{Title: "T", Description: "new", Completed: &completed}
	b, _ := json.Marshal(body)

	mockSvc.On("UpdateTodo", uint(1), &body).Return(&Todo{ID: 1, Title: "T", Description: "new", Completed: true}, nil).Once()

	req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP PUT /todos/1: status=%d resp=%s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
	var resp Todo
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "new", resp.Description)
	assert.True(t, resp.Completed)

	mockSvc.AssertExpectations(t)
}

func TestUpdateTodo_BadRequest_InvalidJSON(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	b := []byte(`{"description":"x"}`) // missing required title
	req := httptest.NewRequest(http.MethodPut, "/todos/1", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP PUT /todos/1 (bad): status=%d resp=%s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTodo_Success(t *testing.T) {
	mockSvc := new(mockTodoService)
	h := NewTodoHandler(mockSvc)
	r := setupRouter(h)

	mockSvc.On("DeleteTodo", uint(3)).Return(nil).Once()

	req := httptest.NewRequest(http.MethodDelete, "/todos/3", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	t.Logf("HTTP DELETE /todos/3: status=%d resp=%s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	assert.Equal(t, "Todo deleted successfully", resp["message"])

	mockSvc.AssertExpectations(t)
}
