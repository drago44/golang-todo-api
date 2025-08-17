package database

import (
	"github.com/drago44/golang-todo-api/internal/config"
	"github.com/drago44/golang-todo-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Init initializes and returns a database connection using the provided config.
func Init(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	// Open a SQLite database named "todo.db"
	db, err := gorm.Open(sqlite.Open(cfg.URL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate runs the database migrations for the Todo model.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Todo{})
}
