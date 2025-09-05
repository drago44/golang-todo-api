package todos

import (
	"errors"
	"fmt"
)

// TodoService defines business logic for managing todos.
type TodoService interface {
	CreateTodo(req *CreateTodoRequest) (*Todo, error)
	GetAllTodos() ([]Todo, error)
	GetTodoByID(id uint) (*Todo, error)
	UpdateTodo(id uint, req *UpdateTodoRequest) (*Todo, error)
	DeleteTodo(id uint) error
}











type todoService struct {
	todoRepo TodoRepository
}

// NewTodoService constructs a TodoService with the provided repository.
func NewTodoService(todoRepo TodoRepository) TodoService {
	return &todoService{todoRepo: todoRepo}
}

// Domain errors returned by TodoService and repository.
var (
	ErrTitleRequired = errors.New("title is required")
	ErrTitleExists   = errors.New("todo with this title already exists")
	ErrNotFound      = errors.New("todo not found")
)

func (s *todoService) CreateTodo(req *CreateTodoRequest) (*Todo, error) {
	// 1. Check if title is required
	if req.Title == "" {
		return nil, ErrTitleRequired
	}

	// 2. Check if title already exists in the database
	exists, err := s.todoRepo.ExistsByTitle(req.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to check title uniqueness: %w", err)
	}
	if exists {
		return nil, ErrTitleExists
	}

	// 3. Create a new Todo
	todo := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	// 4. Save to the database
	if err := s.todoRepo.Create(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) GetAllTodos() ([]Todo, error) {
	return s.todoRepo.GetAll()
}

func (s *todoService) GetTodoByID(id uint) (*Todo, error) {
	return s.todoRepo.GetByID(id)
}

func (s *todoService) UpdateTodo(id uint, req *UpdateTodoRequest) (*Todo, error) {
	// 1. Get the existing Todo
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 2. Track changes
	hasChanges := false

	// 3. Update Title (if provided)
	if req.Title != "" && req.Title != todo.Title {
		// Check if the title is unique
		exists, err := s.todoRepo.ExistsByTitle(req.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to check title uniqueness: %w", err)
		}
		if exists {
			return nil, ErrTitleExists
		}

		todo.Title = req.Title
		hasChanges = true
	}

	// 4. Update Description (if provided)
	if req.Description != "" && req.Description != todo.Description {
		todo.Description = req.Description
		hasChanges = true
	}

	// 5. Update Completed (if provided)
	if req.Completed != nil && *req.Completed != todo.Completed {
		todo.Completed = *req.Completed
		hasChanges = true
	}

	// 6. If there are no changes, return without updating
	if !hasChanges {
		return todo, nil
	}

	// 7. Save the changes
	if err := s.todoRepo.Update(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) DeleteTodo(id uint) error {
	return s.todoRepo.Delete(id)
}
