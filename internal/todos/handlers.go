package todos

import (
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// TodoHandler exposes HTTP handlers for todo resources.
type TodoHandler struct {
	todoService TodoService
}

// NewTodoHandler creates a new TodoHandler instance.
func NewTodoHandler(todoService TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// RegisterTodoRoutes registers todo routes under the provided router group.
func (h *TodoHandler) RegisterTodoRoutes(rg *gin.RouterGroup) {
	todos := rg.Group("/todos")
	{
		todos.POST("", h.CreateTodo)
		todos.GET("", h.GetAllTodos)
		todos.GET("/:id", h.GetTodoByID)
		todos.PUT("/:id", h.UpdateTodo)
		todos.DELETE("/:id", h.DeleteTodo)
	}
}

// sync.Pool removed for simplicity

// ErrorResponse describes an error payload returned by the API.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse describes a simple informational message payload.
type MessageResponse struct {
	Message string `json:"message"`
}

// CreateTodo handles POST /todos and creates a new todo item.
// @Summary Create a new todo
// @Description Create a todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param request body CreateTodoRequest true "Create Todo Request"
// @Success 201 {object} Todo
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *gin.Context) {
	req := new(CreateTodoRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	todo, err := h.todoService.CreateTodo(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTitleRequired):
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		case errors.Is(err, ErrTitleExists):
			c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, todo)
}

// GetAllTodos handles GET /todos and returns all todo items.
// @Summary List todos
// @Description Get all todos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} Todo
// @Failure 500 {object} ErrorResponse
// @Router /todos [get]
func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.todoService.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// GetTodoByID handles GET /todos/{id} to fetch a todo by ID.
// @Summary Get todo by ID
// @Description Get a todo by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} Todo
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	// Parse the ID from the URL parameter and convert it to uint type
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	todo, err := h.todoService.GetTodoByID(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Todo not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, todo)
}

// UpdateTodo handles PUT /todos/{id} to update a todo item by ID.
// @Summary Update todo
// @Description Update a todo by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param request body UpdateTodoRequest true "Update Todo Request"
// @Success 200 {object} Todo
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	req := new(UpdateTodoRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	todo, err := h.todoService.UpdateTodo(uint(id), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Todo not found"})
			return
		case errors.Is(err, ErrTitleExists):
			c.JSON(http.StatusConflict, ErrorResponse{Error: err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo handles DELETE /todos/{id} to remove a todo by ID.
// @Summary Delete todo
// @Description Delete a todo by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	if err := h.todoService.DeleteTodo(uint(id)); err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Todo not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "Todo deleted successfully"})
}
