package todos

import (
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *Todo) error
	GetAll() ([]Todo, error)
	GetByID(id uint) (*Todo, error)
	GetByTitle(title string) (*Todo, error)
	Update(id uint, todo *Todo) error
	Delete(id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(todo *Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) GetAll() ([]Todo, error) {
	var todos []Todo
	err := r.db.Find(&todos).Error
	return todos, err
}

func (r *todoRepository) GetByID(id uint) (*Todo, error) {
	var todo Todo
	err := r.db.First(&todo, id).Error
	return &todo, err
}

func (r *todoRepository) GetByTitle(title string) (*Todo, error) {
	var todo Todo
	err := r.db.Where("title = ?", title).First(&todo).Error
	return &todo, err
}

func (r *todoRepository) Update(id uint, todo *Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&Todo{}, id).Error
}
