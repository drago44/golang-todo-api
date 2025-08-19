package router

import (
	"github.com/drago44/golang-todo-api/internal/todos"
	"github.com/gin-gonic/gin"
)

// Router is a struct that contains the engine and the todo handler
type Router struct {
	engine      *gin.Engine
	todoHandler *todos.TodoHandler
}

// New creates a new router, receiving the handlers through DI
func New(engine *gin.Engine, todoHandler *todos.TodoHandler) *Router {
	r := &Router{
		engine:      engine,
		todoHandler: todoHandler,
	}
	r.setupRoutes()
	return r
}

// setupRoutes sets up all routes
func (r *Router) setupRoutes() {
	// Simple routes
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "API is running"})
	})

	// API v1 group
	v1 := r.engine.Group("/api/v1")

	// Register Todo routes through the injected handler
	r.todoHandler.RegisterTodoRoutes(v1)
}

// GetEngine returns the *gin.Engine for running the server
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
