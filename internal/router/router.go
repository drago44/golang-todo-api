package router

import (
	"github.com/drago44/golang-todo-api/internal/container"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	c      *container.Container
}

func New(engine *gin.Engine, c *container.Container) *Router {
	r := &Router{engine: engine, c: c}
	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "API is running"})
	})

	v1 := r.engine.Group("/api/v1")

	RegisterTodoRoutes(v1, r.c.TodoHandler)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
