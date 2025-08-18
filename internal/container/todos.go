package container

import (
	"github.com/drago44/golang-todo-api/internal/handlers"
	"github.com/drago44/golang-todo-api/internal/repository"
	"github.com/drago44/golang-todo-api/internal/service"
	"gorm.io/gorm"
)

func NewTodoHandler(db *gorm.DB) *handlers.TodoHandler {
	todoRepo := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepo)
	return handlers.NewTodoHandler(todoService)
}
