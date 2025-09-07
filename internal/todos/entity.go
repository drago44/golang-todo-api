package todos

import (
	"time"

	"gorm.io/gorm"
)

// Todo represents a todo item stored in the database.
type Todo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"type:text;uniqueIndex:idx_todos_title_not_deleted,where:deleted_at IS NULL;not null"`
	Description string         `json:"description"`
	Completed   bool           `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
