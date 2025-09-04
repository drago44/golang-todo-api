package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Add performant SQLite options
	dsn = ensureSQLitePragmas(dsn)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		CreateBatchSize:        1000,
	})
	if err != nil {
		return nil, fmt.Errorf("opening sqlite at %s: %w", dsn, err)
	}

	if sqlDB, err2 := db.DB(); err2 == nil {
		sqlDB.SetMaxOpenConns(4)
		sqlDB.SetMaxIdleConns(4)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	}

	return db, nil
}

// Migrate runs the database migrations for all models.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&todos.Todo{})
}

// ensureSQLitePragmas appends performance-friendly PRAGMA options to DSN
func ensureSQLitePragmas(dsn string) string {
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}

	addOpt := func(s, key, pair string) (string, string) {
		if strings.Contains(strings.ToLower(s), strings.ToLower(key+"=")) {
			return s, sep
		}

		if sep == "?" {
			s += "?" + pair
		} else {
			s += "&" + pair
		}

		return s, "&"
	}
	out := dsn
	out, sep = addOpt(out, "_journal_mode", "_journal_mode=WAL")
	out, sep = addOpt(out, "_synchronous", "_synchronous=NORMAL")
	out, sep = addOpt(out, "_busy_timeout", "_busy_timeout=5000")
	out, sep = addOpt(out, "_cache_size", "_cache_size=-20000")
	out, _ = addOpt(out, "_foreign_keys", "_foreign_keys=ON")

	return out
}
