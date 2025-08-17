package repository

import (
	"github.com/drago44/golang-todo-api/internal/models"
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *models.Todo) error
	GetAll() ([]models.Todo, error)
	GetByID(id uint) (*models.Todo, error)
	GetByTitle(title string) (*models.Todo, error)
	Update(id uint, todo *models.Todo) error
	Delete(id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) GetAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := r.db.Find(&todos).Error
	return todos, err
}

func (r *todoRepository) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	return &todo, err
}

func (r *todoRepository) GetByTitle(title string) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Where("title = ?", title).First(&todo).Error
	return &todo, err
}

func (r *todoRepository) Update(id uint, todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}
