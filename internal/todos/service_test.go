package todos

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of TodoRepository for service unit tests
type mockTodoRepository struct{ mock.Mock }

func (m *mockTodoRepository) Create(todo *Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *mockTodoRepository) GetAll() ([]Todo, error) {
	args := m.Called()
	return args.Get(0).([]Todo), args.Error(1)
}

func (m *mockTodoRepository) GetByID(id uint) (*Todo, error) {
	args := m.Called(id)
	if v := args.Get(0); v != nil {
		return v.(*Todo), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockTodoRepository) ExistsByTitle(title string) (bool, error) {
	args := m.Called(title)
	return args.Bool(0), args.Error(1)
}

func (m *mockTodoRepository) Update(todo *Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *mockTodoRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateTodo_Success(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	req := &CreateTodoRequest{Title: "Test", Description: "desc"}
	t.Logf("CreateTodo: preparing request: %+v", req)

	// Simulate title does not exist
	mockRepo.On("ExistsByTitle", "Test").Return(false, nil).Once()
	// Expect create to be called
	mockRepo.On("Create", mock.MatchedBy(func(todo *Todo) bool {
		return todo.Title == "Test" && todo.Description == "desc" && todo.Completed == false
	})).Return(nil).Once()

	created, err := service.CreateTodo(req)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "Test", created.Title)
	assert.Equal(t, "desc", created.Description)
	assert.False(t, created.Completed)
	t.Logf("CreateTodo: created todo: %+v", created)

	mockRepo.AssertExpectations(t)
}

func TestCreateTodo_EmptyTitle(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	_, err := service.CreateTodo(&CreateTodoRequest{Title: "", Description: "x"})
	assert.Error(t, err)
	assert.Equal(t, "title is required", err.Error())
	t.Log("CreateTodo: got expected validation error for empty title")
}

func TestCreateTodo_TitleExists(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	mockRepo.On("ExistsByTitle", "Dup").Return(true, nil).Once()

	_, err := service.CreateTodo(&CreateTodoRequest{Title: "Dup"})
	assert.Error(t, err)
	assert.Equal(t, "todo with this title already exists", err.Error())
	t.Log("CreateTodo: got expected duplicate title error")

	mockRepo.AssertExpectations(t)
}

func TestCreateTodo_ExistsByTitleDbError(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	mockRepo.On("ExistsByTitle", "X").Return(false, errors.New("db down")).Once()

	_, err := service.CreateTodo(&CreateTodoRequest{Title: "X"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check title uniqueness")
	t.Logf("CreateTodo: got expected repo error: %v", err)

	mockRepo.AssertExpectations(t)
}

func TestGetAllTodos(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	expected := []Todo{{ID: 1, Title: "A"}}
	mockRepo.On("GetAll").Return(expected, nil).Once()

	got, err := service.GetAllTodos()
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	t.Logf("GetAllTodos: fetched %d todos", len(got))

	mockRepo.AssertExpectations(t)
}

func TestGetTodoByID(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	mockRepo.On("GetByID", uint(7)).Return(&Todo{ID: 7, Title: "Z"}, nil).Once()

	got, err := service.GetTodoByID(7)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, uint(7), got.ID)
	t.Logf("GetTodoByID: got todo: %+v", got)

	mockRepo.AssertExpectations(t)
}

func TestUpdateTodo_NoChanges(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	existing := &Todo{ID: 3, Title: "A", Description: "d", Completed: false}
	mockRepo.On("GetByID", uint(3)).Return(existing, nil).Once()

	req := &UpdateTodoRequest{}
	got, err := service.UpdateTodo(3, req)
	assert.NoError(t, err)
	assert.Equal(t, existing, got)
	t.Log("UpdateTodo: no changes applied as expected")

	// Ensure Update not called
	mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTodo_TitleConflict(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	mockRepo.On("GetByID", uint(5)).Return(&Todo{ID: 5, Title: "Old"}, nil).Once()
	mockRepo.On("ExistsByTitle", "New").Return(true, nil).Once()

	_, err := service.UpdateTodo(5, &UpdateTodoRequest{Title: "New"})
	assert.Error(t, err)
	assert.Equal(t, "todo with this title already exists", err.Error())
	t.Log("UpdateTodo: got expected title conflict error")

	mockRepo.AssertExpectations(t)
}

func TestUpdateTodo_Success(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	mockRepo.On("GetByID", uint(9)).Return(&Todo{ID: 9, Title: "T", Description: "old", Completed: false}, nil).Once()

	completed := true
	req := &UpdateTodoRequest{Description: "new", Completed: &completed}

	mockRepo.On("Update", mock.MatchedBy(func(todo *Todo) bool {
		return todo.Description == "new" && todo.Completed == true && todo.Title == "T"
	})).Return(nil).Once()

	updated, err := service.UpdateTodo(9, req)
	assert.NoError(t, err)
	assert.Equal(t, "new", updated.Description)
	assert.True(t, updated.Completed)
	t.Logf("UpdateTodo: updated todo: %+v", updated)

	mockRepo.AssertExpectations(t)
}

func TestDeleteTodo(t *testing.T) {
	mockRepo := new(mockTodoRepository)
	service := NewTodoService(mockRepo)

	mockRepo.On("Delete", uint(11)).Return(nil).Once()

	err := service.DeleteTodo(11)
	assert.NoError(t, err)
	t.Log("DeleteTodo: delete returned no error")

	mockRepo.AssertExpectations(t)
}
