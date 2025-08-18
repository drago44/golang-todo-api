package container

import (
	"github.com/drago44/golang-todo-api/internal/handlers"
	"gorm.io/gorm"
)

type Container struct {
	TodoHandler *handlers.TodoHandler
}

func NewContainer(db *gorm.DB) *Container {
	return &Container{
		TodoHandler: NewTodoHandler(db),
	}
}
