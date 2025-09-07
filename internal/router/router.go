// Package router wires HTTP routes for the API.
package router

import (
	"net/http"

	"github.com/drago44/golang-todo-api/internal/todos"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router contains the Gin engine and handlers configuration.
type Router struct {
	engine         *gin.Engine
	todoHandler    *todos.TodoHandler
	swaggerEnabled bool
}

// New creates a new Router and sets up routes.
func New(engine *gin.Engine, todoHandler *todos.TodoHandler, swaggerEnabled bool) *Router {
	r := &Router{
		engine:         engine,
		todoHandler:    todoHandler,
		swaggerEnabled: swaggerEnabled,
	}
	r.setupRoutes()
	return r
}

// setupRoutes sets up all routes
func (r *Router) setupRoutes() {
	// Simple routes
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "API is running"})
	})

	if r.swaggerEnabled {
		// Serve Swagger UI using generated docs package
		r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 group
	v1 := r.engine.Group("/api/v1")

	// Register Todo routes through the injected handler
	r.todoHandler.RegisterTodoRoutes(v1)
}

// GetEngine returns the *gin.Engine for running the server
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
