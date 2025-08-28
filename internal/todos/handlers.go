package todos

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	todoService TodoService
}

func NewTodoHandler(todoService TodoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

func (h *TodoHandler) RegisterTodoRoutes(rg *gin.RouterGroup) {
	todos := rg.Group("/todos")
	{
		todos.POST("/", h.CreateTodo)
		todos.GET("/", h.GetAllTodos)
		todos.GET("/:id", h.GetTodoByID)
		todos.PUT("/:id", h.UpdateTodo)
		todos.DELETE("/:id", h.DeleteTodo)
	}
}

var (
	createTodoReqPool = sync.Pool{New: func() any { return new(CreateTodoRequest) }}
	updateTodoReqPool = sync.Pool{New: func() any { return new(UpdateTodoRequest) }}
)

type errorResponse struct {
	Error string `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	req := createTodoReqPool.Get().(*CreateTodoRequest)
	defer func() {
		*req = CreateTodoRequest{}
		createTodoReqPool.Put(req)
	}()
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	todo, err := h.todoService.CreateTodo(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) GetAllTodos(c *gin.Context) {
	todos, err := h.todoService.GetAllTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	// Parse the ID from the URL parameter and convert it to uint type
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Invalid ID"})
		return
	}

	todo, err := h.todoService.GetTodoByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse{Error: "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Invalid ID"})
		return
	}

	req := updateTodoReqPool.Get().(*UpdateTodoRequest)
	defer func() {
		*req = UpdateTodoRequest{}
		updateTodoReqPool.Put(req)
	}()
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	todo, err := h.todoService.UpdateTodo(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Invalid ID"})
		return
	}

	if err := h.todoService.DeleteTodo(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, messageResponse{Message: "Todo deleted successfully"})
}
