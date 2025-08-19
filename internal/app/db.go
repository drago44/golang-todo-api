package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drago44/golang-todo-api/internal/todos"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Init initializes and returns a database connection using the provided config.
func Init(cfg *DatabaseConfig) (*gorm.DB, error) {
	dsn := strings.TrimSpace(cfg.URL)
	if strings.Contains(dsn, "/") || strings.Contains(dsn, string(os.PathSeparator)) {
		dir := filepath.Dir(dsn)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, fmt.Errorf("failed to create db directory %s: %w", dir, err)
			}
		}
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("opening sqlite at %s: %w", dsn, err)
	}
	return db, nil
}

// Migrate runs the database migrations for all models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&todos.Todo{})
}
