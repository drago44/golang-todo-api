package todos

import (
	"errors"

	"gorm.io/gorm"
)

// TodoRepository defines persistence operations for Todo entities.
type TodoRepository interface {
	Create(todo *Todo) error
	GetAll() ([]Todo, error)
	GetByID(id uint) (*Todo, error)
	ExistsByTitle(title string) (bool, error)
	Update(todo *Todo) error
	Delete(id uint) error
}

type todoRepository struct {
	db *gorm.DB
}

// NewTodoRepository creates a GORM-backed TodoRepository.
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) ExistsByTitle(title string) (bool, error) {
	var todo Todo
	res := r.db.Model(&Todo{}).
		Select("id").
		Where("title = ?", title).
		Limit(1).
		Take(&todo)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, res.Error
	}
	return true, nil
}

func (r *todoRepository) Update(todo *Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	var todo Todo
	res := r.db.Delete(&todo, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
