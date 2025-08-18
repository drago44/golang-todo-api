package router

import (
	"github.com/drago44/golang-todo-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterTodoRoutes(rg *gin.RouterGroup, todoHandler *handlers.TodoHandler) {
	todos := rg.Group("/todos")
	{
		todos.POST("/", todoHandler.CreateTodo)
		todos.GET("/", todoHandler.GetAllTodos)
		todos.GET("/:id", todoHandler.GetTodoByID)
		todos.PUT("/:id", todoHandler.UpdateTodo)
		todos.DELETE("/:id", todoHandler.DeleteTodo)
	}
}
