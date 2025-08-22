package todos

import (
	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *Todo) error
	GetAll() ([]Todo, error)
	GetByID(id uint) (*Todo, error)
	ExistsByTitle(title string) (bool, error)
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

func (r *todoRepository) ExistsByTitle(title string) (bool, error) {
	var count int64
	if err := r.db.Model(&Todo{}).Where("title = ?", title).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *todoRepository) Update(id uint, todo *Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&Todo{}, id).Error
}
