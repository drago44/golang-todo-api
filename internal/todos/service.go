package todos

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

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

func NewTodoService(todoRepo TodoRepository) TodoService {
	return &todoService{todoRepo: todoRepo}
}

func (s *todoService) CreateTodo(req *CreateTodoRequest) (*Todo, error) {
	// 1. Check if title is required
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	// 2. Check if title already exists in the database
	existingTodo, err := s.todoRepo.GetByTitle(req.Title)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// If the error is not "record not found", then it's a database error
		return nil, fmt.Errorf("failed to check title uniqueness: %w", err)
	}

	if existingTodo != nil {
		return nil, errors.New("todo with this title already exists")
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
		existingTodo, err := s.todoRepo.GetByTitle(req.Title)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check title uniqueness: %w", err)
		}
		if existingTodo != nil && existingTodo.ID != id {
			return nil, errors.New("todo with this title already exists")
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
	if err := s.todoRepo.Update(id, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) DeleteTodo(id uint) error {
	return s.todoRepo.Delete(id)
}
